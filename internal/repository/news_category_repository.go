package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsCategoryRepository interface {
	FindActiveBySlug(ctx context.Context, slug string, lang string) (*entity.NewsCategory, error)
	FindAll(ctx context.Context) ([]entity.NewsCategory, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.NewsCategory, error)
	Create(ctx context.Context, item *entity.NewsCategory) error
	Update(ctx context.Context, item *entity.NewsCategory) error
	UpsertTranslation(ctx context.Context, item *entity.NewsCategoryTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type newsCategoryRepositoryImpl struct {
	db *gorm.DB
}

func NewNewsCategoryRepository(db *gorm.DB) NewsCategoryRepository {
	return &newsCategoryRepositoryImpl{
		db: db,
	}
}

func (r *newsCategoryRepositoryImpl) FindActiveBySlug(
	ctx context.Context,
	slug string,
	lang string,
) (*entity.NewsCategory, error) {

	var category entity.NewsCategory

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("id IN (?)", r.subQueryCategoryIDBySlug(slug, lang)).
		First(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *newsCategoryRepositoryImpl) FindAll(ctx context.Context) ([]entity.NewsCategory, error) {
	var categories []entity.NewsCategory
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("created_at DESC").
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *newsCategoryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.NewsCategory, error) {
	var category entity.NewsCategory
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *newsCategoryRepositoryImpl) Create(ctx context.Context, item *entity.NewsCategory) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *newsCategoryRepositoryImpl) Update(ctx context.Context, item *entity.NewsCategory) error {
	return r.db.WithContext(ctx).
		Model(&entity.NewsCategory{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"is_active": item.IsActive,
		}).Error
}

func (r *newsCategoryRepositoryImpl) UpsertTranslation(ctx context.Context, item *entity.NewsCategoryTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "category_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "slug", "description", "updated_at"}),
		}).
		Create(item).Error
}

func (r *newsCategoryRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.NewsCategory{}, "id = ?", id).Error
}

func (r *newsCategoryRepositoryImpl) subQueryCategoryIDBySlug(
	slug string,
	lang string,
) *gorm.DB {

	return r.db.
		Table("news_category_translations").
		Select("category_id").
		Where("slug = ?", slug).
		Where("language = ?", lang)
}
