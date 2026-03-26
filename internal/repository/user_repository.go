package repository

import (
	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByEmailOrUsername(identifier string) (*entity.User, error)
	Update(user *entity.User) error
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) FindByID(id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("id = ?", id).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).Take(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByEmailOrUsername(identifier string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) Update(user *entity.User) error {
	return r.db.Save(user).Error
}
