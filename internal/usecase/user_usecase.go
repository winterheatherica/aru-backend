package usecase

import (
	"aru-backend/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate
	UserRepo repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:       db,
		Log:      logger,
		Validate: validate,
		UserRepo: userRepo,
	}
}
