package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type HeroRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.HeroSlide, error)
}

type heroRepositoryImpl struct {
	db *gorm.DB
}

func NewHeroRepository(db *gorm.DB) HeroRepository {
	return &heroRepositoryImpl{db: db}
}

func (r *heroRepositoryImpl) FindActiveByLanguage(ctx context.Context, lang string) ([]entity.HeroSlide, error) {
	var slides []entity.HeroSlide

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Order("order_index ASC").
		Find(&slides).Error

	if err != nil {
		return nil, err
	}

	return slides, nil
}
