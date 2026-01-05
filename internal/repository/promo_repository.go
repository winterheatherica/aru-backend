package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type PromoRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.PromoSlide, error)
}

type promoRepositoryImpl struct {
	db *gorm.DB
}

func NewPromoRepository(db *gorm.DB) PromoRepository {
	return &promoRepositoryImpl{db: db}
}

func (r *promoRepositoryImpl) FindActiveByLanguage(
	ctx context.Context,
	lang string,
) ([]entity.PromoSlide, error) {
	var slides []entity.PromoSlide

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
