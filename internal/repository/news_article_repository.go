package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type NewsArticleRepository interface {
	FindActiveBySlug(ctx context.Context, slug string, lang string) (*entity.NewsArticle, error)
}

type newsArticleRepositoryImpl struct {
	db *gorm.DB
}

func NewNewsArticleRepository(db *gorm.DB) NewsArticleRepository {
	return &newsArticleRepositoryImpl{
		db: db,
	}
}

func (r *newsArticleRepositoryImpl) FindActiveBySlug(
	ctx context.Context,
	slug string,
	lang string,
) (*entity.NewsArticle, error) {

	var article entity.NewsArticle

	subQuery := r.subQueryArticleIDBySlug(slug, lang)

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("id IN (?)", subQuery).
		Preload("Translations", "language = ?", lang).
		Preload("Categories.Translations", "language = ?", lang).
		First(&article).Error

	if err != nil {
		return nil, err
	}

	return &article, nil
}

func (r *newsArticleRepositoryImpl) subQueryArticleIDBySlug(
	slug string,
	lang string,
) *gorm.DB {

	return r.db.
		Table("news_article_translations").
		Select("article_id").
		Where("slug = ?", slug).
		Where("language = ?", lang)
}
