package repository

import (
	"aru-backend/internal/entity"
	"context"

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
		Joins("JOIN service_pricing_tier_translations t ON t.tier_id = service_pricing_tiers.id").
		Where("service_pricing_tiers.is_active = ?", true).
		Where("service_pricing_tiers.service = ?", service).
		Where("t.language = ?", lang).
		Order("service_pricing_tiers.created_at DESC").
		Limit(3).
		Preload("Translations", "language = ?", lang).
		Find(&tiers).Error

	if err != nil {
		return nil, err
	}

	return tiers, nil
}
