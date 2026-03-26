package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AwardRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.Award, error)
	FindAll(ctx context.Context) ([]entity.Award, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Award, error)
	Create(ctx context.Context, item *entity.Award) error
	Update(ctx context.Context, item *entity.Award) error
	UpsertTranslation(ctx context.Context, tr *entity.AwardTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *awardRepositoryImpl) FindAll(ctx context.Context) ([]entity.Award, error) {
	var items []entity.Award
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active DESC, year DESC, order_index ASC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *awardRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Award, error) {
	var item entity.Award
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *awardRepositoryImpl) Create(ctx context.Context, item *entity.Award) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *awardRepositoryImpl) Update(ctx context.Context, item *entity.Award) error {
	return r.db.WithContext(ctx).
		Model(&entity.Award{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"image_path":  item.ImagePath,
			"year":        item.Year,
			"order_index": item.OrderIndex,
			"is_active":   item.IsActive,
			"uploaded_by": item.UploadedBy,
		}).Error
}

func (r *awardRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.AwardTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "award_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"alt":         tr.Alt,
				"title":       tr.Title,
				"label":       tr.Label,
				"description": tr.Description,
			}),
		}).
		Create(tr).Error
}

func (r *awardRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Award{}, "id = ?", id).Error
}
