package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type AwardRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.Award, error)
}

type awardRepositoryImpl struct {
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) AwardRepository {
	return &awardRepositoryImpl{
		db: db,
	}
}

func (r *awardRepositoryImpl) FindActiveByLanguage(
	ctx context.Context,
	lang string,
) ([]entity.Award, error) {

	var awards []entity.Award

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Order("year ASC, order_index ASC").
		Find(&awards).Error

	if err != nil {
		return nil, err
	}

	return awards, nil
}
