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
		Joins("JOIN news_category_translations t ON t.category_id = news_categories.id").
		Where("news_categories.is_active = ?", true).
		Where("t.language = ?", lang).
		Where("t.slug = ?", slug).
		Preload("Translations", "language = ?", lang).
		First(&category).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}
