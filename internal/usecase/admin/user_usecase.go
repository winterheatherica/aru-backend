package admin

import (
	"context"
	"fmt"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	List(ctx context.Context) ([]model.UserAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserAdminItem, error)
	Create(ctx context.Context, input model.UserAdminUpsertInput) (*model.UserAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.UserAdminUpsertInput) (*model.UserAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type userUsecaseImpl struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{repo: repo}
}

func (u *userUsecaseImpl) List(ctx context.Context) ([]model.UserAdminItem, error) {
	_ = ctx
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	out := make([]model.UserAdminItem, 0, len(users))
	for _, it := range users {
		out = append(out, toUserAdminItem(it))
	}
	return out, nil
}

func (u *userUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.UserAdminItem, error) {
	_ = ctx
	it, err := u.repo.FindByUUID(id)
	if err != nil {
		return nil, err
	}
	res := toUserAdminItem(*it)
	return &res, nil
}

func (u *userUsecaseImpl) Create(ctx context.Context, input model.UserAdminUpsertInput) (*model.UserAdminItem, error) {
	_ = ctx
	if err := validateUserInput(input, true); err != nil {
		return nil, err
	}
	if _, err := u.repo.FindByEmail(strings.TrimSpace(input.Email)); err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	active := true
	if input.Active != nil {
		active = *input.Active
	}
	emailVerified := true
	if input.EmailVerified != nil {
		emailVerified = *input.EmailVerified
	}

	username := strings.TrimSpace(input.Username)
	fullName := strings.TrimSpace(input.FullName)
	user := &entity.User{
		Email:         strings.TrimSpace(input.Email),
		Username:      &username,
		FullName:      &fullName,
		Password:      string(hash),
		Role:          normalizeRole(input.Role),
		IsActive:      active,
		EmailVerified: emailVerified,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	res := toUserAdminItem(*user)
	return &res, nil
}

func (u *userUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.UserAdminUpsertInput) (*model.UserAdminItem, error) {
	_ = ctx
	if err := validateUserInput(input, false); err != nil {
		return nil, err
	}

	user, err := u.repo.FindByUUID(id)
	if err != nil {
		return nil, err
	}

	if email := strings.TrimSpace(input.Email); email != "" {
		user.Email = email
	}
	if username := strings.TrimSpace(input.Username); username != "" {
		user.Username = &username
	}
	if fullName := strings.TrimSpace(input.FullName); fullName != "" {
		user.FullName = &fullName
	}
	if role := strings.TrimSpace(input.Role); role != "" {
		user.Role = normalizeRole(role)
	}
	if input.Active != nil {
		user.IsActive = *input.Active
	}
	if input.EmailVerified != nil {
		user.EmailVerified = *input.EmailVerified
	}
	if strings.TrimSpace(input.Password) != "" {
		hash, hErr := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if hErr != nil {
			return nil, hErr
		}
		user.Password = string(hash)
	}

	if err := u.repo.Update(user); err != nil {
		return nil, err
	}

	updated, err := u.repo.FindByUUID(id)
	if err != nil {
		res := toUserAdminItem(*user)
		return &res, nil
	}
	res := toUserAdminItem(*updated)
	return &res, nil
}

func (u *userUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	_ = ctx
	_, err := u.repo.FindByUUID(id)
	if err != nil {
		return err
	}
	return u.repo.DeleteByUUID(id)
}

func validateUserInput(input model.UserAdminUpsertInput, requirePassword bool) error {
	if strings.TrimSpace(input.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(input.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(input.FullName) == "" {
		return fmt.Errorf("full_name is required")
	}
	if requirePassword && strings.TrimSpace(input.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if role := strings.TrimSpace(input.Role); role != "" {
		role = strings.ToUpper(role)
		if role != "ADMIN" && role != "EDITOR" && role != "MEMBER" {
			return fmt.Errorf("role must be ADMIN, EDITOR, or MEMBER")
		}
	}
	return nil
}

func normalizeRole(role string) string {
	up := strings.ToUpper(strings.TrimSpace(role))
	if up == "" {
		return "MEMBER"
	}
	return up
}

func ptrValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func toUserAdminItem(user entity.User) model.UserAdminItem {
	return model.UserAdminItem{
		ID:            user.ID,
		Email:         user.Email,
		Username:      ptrValue(user.Username),
		FullName:      ptrValue(user.FullName),
		Role:          user.Role,
		Active:        user.IsActive,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
