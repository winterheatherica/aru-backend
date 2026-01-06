package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type NewsArticleRepository interface {
	FindActiveBySlug(ctx context.Context, slug string, lang string) (*entity.NewsArticle, error)
	FindActiveCardList(ctx context.Context, lang string, year *int, limit int, offset int) ([]entity.NewsArticle, error)
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

func (r *newsArticleRepositoryImpl) FindActiveCardList(ctx context.Context, lang string, year *int, limit int, offset int) ([]entity.NewsArticle, error) {
	var articles []entity.NewsArticle

	q := r.db.WithContext(ctx).
		Model(&entity.NewsArticle{}).
		Joins("JOIN news_article_translations t ON t.article_id = news_articles.id").
		Where("news_articles.is_active = ?", true).
		Where("news_articles.deleted_at IS NULL").
		Where("t.language = ?", lang).
		Order("news_articles.published_at DESC").
		Limit(limit).
		Offset(offset).
		Preload("Translations", "language = ?", lang)

	if year != nil {
		q = q.Where("EXTRACT(YEAR FROM news_articles.published_at) = ?", *year)
	}

	err := q.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
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
