package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type NewsArticleRepository interface {
	FindActiveByID(ctx context.Context, id string, lang string) (*entity.NewsArticle, error)
	FindActiveCardList(ctx context.Context, lang string, year *int, limit int, offset int) ([]entity.NewsArticle, error)
	FindLatest(ctx context.Context, lang string, limit int) ([]entity.NewsArticle, error)
	FindActiveYears(ctx context.Context) ([]int, error)
}

type newsArticleRepositoryImpl struct {
	db *gorm.DB
}

func NewNewsArticleRepository(db *gorm.DB) NewsArticleRepository {
	return &newsArticleRepositoryImpl{
		db: db,
	}
}

func (r *newsArticleRepositoryImpl) FindActiveByID(
	ctx context.Context,
	id string,
	lang string,
) (*entity.NewsArticle, error) {

	var article entity.NewsArticle

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("id = ?", id).
		Preload("Translations", "language = ?", lang).
		Preload("Categories.Translations", "language = ?", lang).
		First(&article).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &article, nil
}

func (r *newsArticleRepositoryImpl) FindActiveCardList(
	ctx context.Context,
	lang string,
	year *int,
	limit int,
	offset int,
) ([]entity.NewsArticle, error) {

	var articles []entity.NewsArticle

	q := r.db.WithContext(ctx).
		Model(&entity.NewsArticle{}).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Order("published_at DESC").
		Limit(limit).
		Offset(offset)

	if year != nil {
		q = q.Where("EXTRACT(YEAR FROM published_at) = ?", *year)
	}

	err := q.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *newsArticleRepositoryImpl) FindLatest(
	ctx context.Context,
	lang string,
	limit int,
) ([]entity.NewsArticle, error) {

	var articles []entity.NewsArticle

	err := r.db.WithContext(ctx).
		Model(&entity.NewsArticle{}).
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Preload("Translations", "language = ?", lang).
		Order("published_at DESC").
		Limit(limit).
		Find(&articles).Error

	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *newsArticleRepositoryImpl) FindActiveYears(
	ctx context.Context,
) ([]int, error) {

	var years []int

	err := r.db.WithContext(ctx).
		Table("news_articles").
		Select("DISTINCT EXTRACT(YEAR FROM published_at) AS year").
		Where("is_active = ?", true).
		Where("deleted_at IS NULL").
		Order("year DESC").
		Pluck("year", &years).Error

	if err != nil {
		return nil, err
	}

	return years, nil
}
