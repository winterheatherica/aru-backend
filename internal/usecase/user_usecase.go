package usecase

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB        *gorm.DB
	Log       *logrus.Logger
	Validate  *validator.Validate
	UserRepo  repository.UserRepository
	SessRepo  repository.SessionRepository
	JWTSecret string
	JWTExpiry time.Duration
}

func NewUserUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	userRepo repository.UserRepository,
	sessRepo repository.SessionRepository,
	jwtSecret string,
	jwtExpiry time.Duration,
) *UserUseCase {
	return &UserUseCase{
		DB:        db,
		Log:       logger,
		Validate:  validate,
		UserRepo:  userRepo,
		SessRepo:  sessRepo,
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
	}
}

func verifyPassword(input, stored string) bool {
	if stored == "" {
		return false
	}
	if bcrypt.CompareHashAndPassword([]byte(stored), []byte(input)) == nil {
		return true
	}
	return input == stored
}

func (u *UserUseCase) signJWT(userID, role, jwtID string, expiresAt time.Time) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"jti":  jwtID,
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.JWTSecret))
}

func (u *UserUseCase) parseJWT(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(u.JWTSecret), nil
	})
	return err
}

func (u *UserUseCase) LoginWithPassword(identifier, password, ip, userAgent string) (*model.AuthResponse, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" || password == "" {
		return nil, errors.New("identifier and password are required")
	}

	user, err := u.UserRepo.FindByEmailOrUsername(identifier)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}

	if !verifyPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.LoginCount += 1
	if ip != "" {
		user.LastLoginIP = &ip
	}
	if userAgent != "" {
		user.LastLoginAgent = &userAgent
	}
	if err := u.UserRepo.Update(user); err != nil {
		return nil, err
	}

	jwtID := uuid.NewString()
	expiresAt := time.Now().Add(u.JWTExpiry)
	accessToken, err := u.signJWT(user.ID.String(), user.Role, jwtID, expiresAt)
	if err != nil {
		return nil, err
	}

	session := &entity.Session{
		UserID:           user.ID,
		AccessToken:      &accessToken,
		RefreshTokenHash: "-",
		JWTID:            jwtID,
		ExpiresAt:        expiresAt,
		IsCurrent:        true,
	}
	if ip != "" {
		session.IPAddress = &ip
	}
	if userAgent != "" {
		session.UserAgent = &userAgent
	}
	if err := u.SessRepo.Create(session); err != nil {
		return nil, err
	}

	return converter.UserToAuthResponse(user, accessToken, ""), nil
}

func (u *UserUseCase) GetCurrentUserByAccessToken(token string) (*model.UserResponse, error) {
	if err := u.parseJWT(token); err != nil {
		return nil, err
	}

	session, err := u.SessRepo.FindActiveByAccessToken(token)
	if err != nil {
		return nil, err
	}
	_ = u.SessRepo.UpdateLastUsed(session.ID.String())
	return converter.UserToResponse(&session.User), nil
}

func (u *UserUseCase) LogoutByAccessToken(token string) error {
	if token == "" {
		return nil
	}
	return u.SessRepo.RevokeByAccessToken(token, "USER_LOGOUT")
}
