package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServicePricingRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServicePricingTier, error)
	FindByService(ctx context.Context, service string) ([]entity.ServicePricingTier, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ServicePricingTier, error)
	Create(ctx context.Context, item *entity.ServicePricingTier) error
	Update(ctx context.Context, item *entity.ServicePricingTier) error
	UpsertTranslation(ctx context.Context, tr *entity.ServicePricingTierTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *servicePricingRepositoryImpl) FindByService(ctx context.Context, service string) ([]entity.ServicePricingTier, error) {
	var items []entity.ServicePricingTier
	err := r.db.WithContext(ctx).
		Model(&entity.ServicePricingTier{}).
		Where("service = ?", service).
		Preload("Translations").
		Order("updated_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *servicePricingRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ServicePricingTier, error) {
	var item entity.ServicePricingTier
	err := r.db.WithContext(ctx).
		Model(&entity.ServicePricingTier{}).
		Where("id = ?", id).
		Preload("Translations").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *servicePricingRepositoryImpl) Create(ctx context.Context, item *entity.ServicePricingTier) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *servicePricingRepositoryImpl) Update(ctx context.Context, item *entity.ServicePricingTier) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *servicePricingRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.ServicePricingTierTranslation) error {
	return r.db.WithContext(ctx).
		Where("tier_id = ? AND language = ?", tr.TierID, tr.Language).
		Assign(map[string]any{
			"name":        tr.Name,
			"description": tr.Description,
			"features":    tr.Features,
		}).
		FirstOrCreate(tr).Error
}

func (r *servicePricingRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ServicePricingTier{}, "id = ?", id).Error
}
