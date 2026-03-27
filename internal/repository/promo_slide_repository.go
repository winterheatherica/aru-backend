package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PromoRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.PromoSlide, error)
	FindAll(ctx context.Context) ([]entity.PromoSlide, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.PromoSlide, error)
	Create(ctx context.Context, item *entity.PromoSlide) error
	Update(ctx context.Context, item *entity.PromoSlide) error
	UpsertTranslation(ctx context.Context, tr *entity.PromoSlideTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *promoRepositoryImpl) FindAll(ctx context.Context) ([]entity.PromoSlide, error) {
	var items []entity.PromoSlide
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active DESC, order_index ASC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *promoRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.PromoSlide, error) {
	var item entity.PromoSlide
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *promoRepositoryImpl) Create(ctx context.Context, item *entity.PromoSlide) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *promoRepositoryImpl) Update(ctx context.Context, item *entity.PromoSlide) error {
	return r.db.WithContext(ctx).
		Model(&entity.PromoSlide{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"image_path":  item.ImagePath,
			"order_index": item.OrderIndex,
			"is_active":   item.IsActive,
			"uploaded_by": item.UploadedBy,
		}).Error
}

func (r *promoRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.PromoSlideTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "promo_slide_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"alt":   tr.Alt,
				"title": tr.Title,
			}),
		}).
		Create(tr).Error
}

func (r *promoRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.PromoSlide{}, "id = ?", id).Error
}
