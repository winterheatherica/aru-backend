package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PartnerRepository interface {
	FindActiveForGrid(ctx context.Context, lang string) ([]entity.Partner, error)
	FindActiveForScroller(ctx context.Context, lang string) ([]entity.Partner, error)
	FindAll(ctx context.Context) ([]entity.Partner, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Partner, error)
	Create(ctx context.Context, item *entity.Partner) error
	Update(ctx context.Context, item *entity.Partner) error
	UpsertTranslation(ctx context.Context, tr *entity.PartnerTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *partnerRepositoryImpl) FindAll(ctx context.Context) ([]entity.Partner, error) {
	var items []entity.Partner
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active_partner_grid DESC, is_active_partner_scroller DESC, order_index ASC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *partnerRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Partner, error) {
	var item entity.Partner
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *partnerRepositoryImpl) Create(ctx context.Context, item *entity.Partner) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *partnerRepositoryImpl) Update(ctx context.Context, item *entity.Partner) error {
	return r.db.WithContext(ctx).
		Model(&entity.Partner{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"image_path":                 item.ImagePath,
			"order_index":                item.OrderIndex,
			"is_active_partner_grid":     item.IsActivePartnerGrid,
			"is_active_partner_scroller": item.IsActivePartnerScroller,
			"uploaded_by":                item.UploadedBy,
		}).Error
}

func (r *partnerRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.PartnerTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "partner_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"alt":         tr.Alt,
				"title":       tr.Title,
				"description": tr.Description,
			}),
		}).
		Create(tr).Error
}

func (r *partnerRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Partner{}, "id = ?", id).Error
}
