package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type NewsCategoryRepository interface {
	FindActiveBySlug(ctx context.Context, slug string, lang string) (*entity.NewsCategory, error)
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
