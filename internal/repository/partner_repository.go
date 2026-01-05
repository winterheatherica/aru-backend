package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type PartnerRepository interface {
	FindActiveForGrid(ctx context.Context, lang string) ([]entity.Partner, error)
	FindActiveForScroller(ctx context.Context, lang string) ([]entity.Partner, error)
}

type partnerRepositoryImpl struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepositoryImpl{
		db: db,
	}
}

func (r *partnerRepositoryImpl) FindActiveForGrid(
	ctx context.Context,
	lang string,
) ([]entity.Partner, error) {

	var partners []entity.Partner

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active_partner_grid = ?", true).
		Order("order_index ASC").
		Find(&partners).Error

	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (r *partnerRepositoryImpl) FindActiveForScroller(
	ctx context.Context,
	lang string,
) ([]entity.Partner, error) {

	var partners []entity.Partner

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active_partner_scroller = ?", true).
		Order("order_index ASC").
		Find(&partners).Error

	if err != nil {
		return nil, err
	}

	return partners, nil
}
