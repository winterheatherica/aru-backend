package repository

import "aru-backend/internal/entity"

type UserRepository interface {
	Create(*entity.User) error
	FindByID(id string) (*entity.User, error)
	Update(*entity.User) error
}
