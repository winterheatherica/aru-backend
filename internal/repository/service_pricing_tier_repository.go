package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type ServicePricingRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServicePricingTier, error)
}

type servicePricingRepositoryImpl struct {
	db *gorm.DB
}

func NewServicePricingRepository(db *gorm.DB) ServicePricingRepository {
	return &servicePricingRepositoryImpl{db: db}
}

func (r *servicePricingRepositoryImpl) FindActiveByService(
	ctx context.Context,
	service string,
	lang string,
) ([]entity.ServicePricingTier, error) {

	var tiers []entity.ServicePricingTier

	err := r.db.WithContext(ctx).
		Model(&entity.ServicePricingTier{}).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("service = ?", service).
		Order("created_at DESC").
		Limit(3).
		Find(&tiers).Error

	if err != nil {
		return nil, err
	}

	return tiers, nil
}
