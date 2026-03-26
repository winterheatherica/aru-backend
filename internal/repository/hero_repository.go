package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HeroRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.HeroSlide, error)
	FindAll(ctx context.Context) ([]entity.HeroSlide, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.HeroSlide, error)
	Create(ctx context.Context, slide *entity.HeroSlide) error
	Update(ctx context.Context, slide *entity.HeroSlide) error
	UpsertTranslation(ctx context.Context, tr *entity.HeroSlideTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *heroRepositoryImpl) FindAll(ctx context.Context) ([]entity.HeroSlide, error) {
	var slides []entity.HeroSlide
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active DESC, order_index ASC, created_at DESC").
		Find(&slides).Error; err != nil {
		return nil, err
	}
	return slides, nil
}

func (r *heroRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.HeroSlide, error) {
	var slide entity.HeroSlide
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&slide, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &slide, nil
}

func (r *heroRepositoryImpl) Create(ctx context.Context, slide *entity.HeroSlide) error {
	return r.db.WithContext(ctx).Create(slide).Error
}

func (r *heroRepositoryImpl) Update(ctx context.Context, slide *entity.HeroSlide) error {
	return r.db.WithContext(ctx).
		Model(&entity.HeroSlide{}).
		Where("id = ?", slide.ID).
		Updates(map[string]any{
			"image_path":  slide.ImagePath,
			"order_index": slide.OrderIndex,
			"is_active":   slide.IsActive,
			"banner":      slide.Banner,
			"uploaded_by": slide.UploadedBy,
		}).Error
}

func (r *heroRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.HeroSlideTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "hero_slide_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"alt":         tr.Alt,
				"title":       tr.Title,
				"description": tr.Description,
			}),
		}).
		Create(tr).Error
}

func (r *heroRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.HeroSlide{}, "id = ?", id).Error
}
