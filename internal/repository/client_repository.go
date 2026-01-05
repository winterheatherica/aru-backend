package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type ClientRepository interface {
	FindActiveForScroller(ctx context.Context, lang string) ([]entity.Client, error)
}

type clientRepositoryImpl struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepositoryImpl{
		db: db,
	}
}

func (r *clientRepositoryImpl) FindActiveForScroller(
	ctx context.Context,
	lang string,
) ([]entity.Client, error) {

	var clients []entity.Client

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active_client_scroller = ?", true).
		Order("order_index ASC").
		Find(&clients).Error

	if err != nil {
		return nil, err
	}

	return clients, nil
}
