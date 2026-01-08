package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type HistoryRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.History, error)
}

type historyRepositoryImpl struct {
	db *gorm.DB
}

func NewHistoryRepository(db *gorm.DB) HistoryRepository {
	return &historyRepositoryImpl{
		db: db,
	}
}

func (r *historyRepositoryImpl) FindActiveByLanguage(
	ctx context.Context,
	lang string,
) ([]entity.History, error) {

	var histories []entity.History

	err := r.db.WithContext(ctx).
		Where("language = ?", lang).
		Where("is_active = ?", true).
		Order("year DESC NULLS LAST").
		Find(&histories).Error

	if err != nil {
		return nil, err
	}

	return histories, nil
}
