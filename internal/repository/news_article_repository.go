package repository

import (
	"context"
	"strings"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NewsArticleRepository interface {
	FindActiveByID(ctx context.Context, id string, lang string) (*entity.NewsArticle, error)
	FindPublisherByArticleID(ctx context.Context, id string) (*entity.User, error)
	ResolveIDBySlug(ctx context.Context, slug string) (string, error)
	FindSlugByIDAndLang(ctx context.Context, id, lang string) (string, error)
	FindActiveCardList(ctx context.Context, lang string, year *int, limit int, offset int) ([]entity.NewsArticle, error)
	FindLatest(ctx context.Context, lang string, limit int) ([]entity.NewsArticle, error)
	FindActiveYears(ctx context.Context) ([]int, error)

	FindAll(ctx context.Context) ([]entity.NewsArticle, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.NewsArticle, error)
	Create(ctx context.Context, item *entity.NewsArticle) error
	Update(ctx context.Context, item *entity.NewsArticle) error
	UpsertTranslation(ctx context.Context, item *entity.NewsArticleTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	IsSlugExists(ctx context.Context, slug, lang string, excludeArticleID *uuid.UUID) (bool, error)
	ReplaceCategories(ctx context.Context, articleID uuid.UUID, categoryIDs []uuid.UUID) error
	FindCategoriesByIDs(ctx context.Context, categoryIDs []uuid.UUID) ([]entity.NewsCategory, error)
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

func (r *newsArticleRepositoryImpl) FindPublisherByArticleID(
	ctx context.Context,
	id string,
) (*entity.User, error) {
	var user entity.User

	err := r.db.WithContext(ctx).
		Model(&entity.User{}).
		Joins("JOIN news_articles na ON na.uploaded_by = users.id").
		Where("na.id = ?", id).
		Limit(1).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
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
		Order("published_at DESC")

	if limit > 0 {
		q = q.Limit(limit)
		if offset > 0 {
			q = q.Offset(offset)
		}
	}

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

func (r *newsArticleRepositoryImpl) FindActiveYears(ctx context.Context) ([]int, error) {
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

func (r *newsArticleRepositoryImpl) ResolveIDBySlug(
	ctx context.Context,
	slug string,
) (string, error) {

	var id string

	err := r.db.WithContext(ctx).
		Table("news_article_translations").
		Select("article_id").
		Where("slug = ?", slug).
		Limit(1).
		Scan(&id).Error

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *newsArticleRepositoryImpl) FindSlugByIDAndLang(
	ctx context.Context,
	id, lang string,
) (string, error) {

	var slug string

	err := r.db.WithContext(ctx).
		Table("news_article_translations").
		Select("slug").
		Where("article_id = ?", id).
		Where("language = ?", lang).
		Limit(1).
		Scan(&slug).Error

	if err != nil {
		return "", err
	}

	if slug == "" {
		return "", nil
	}

	return slug, nil
}

func (r *newsArticleRepositoryImpl) FindAll(ctx context.Context) ([]entity.NewsArticle, error) {
	var items []entity.NewsArticle
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Categories.Translations").
		Order("published_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *newsArticleRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.NewsArticle, error) {
	var item entity.NewsArticle
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Categories.Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *newsArticleRepositoryImpl) Create(ctx context.Context, item *entity.NewsArticle) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *newsArticleRepositoryImpl) Update(ctx context.Context, item *entity.NewsArticle) error {
	return r.db.WithContext(ctx).
		Model(&entity.NewsArticle{}).
		Where("id = ?", item.ID).
		Updates(map[string]interface{}{
			"is_active":    item.IsActive,
			"published_at": item.PublishedAt,
			"image_path":   item.ImagePath,
			"uploaded_by":  item.UploadedBy,
		}).Error
}

func (r *newsArticleRepositoryImpl) UpsertTranslation(ctx context.Context, item *entity.NewsArticleTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "article_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"slug", "title", "content", "meta_title", "meta_description", "meta_keywords", "updated_at"}),
		}).
		Create(item).Error
}

func (r *newsArticleRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.NewsArticle{}, "id = ?", id).Error
}

func (r *newsArticleRepositoryImpl) IsSlugExists(ctx context.Context, slug, lang string, excludeArticleID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).
		Table("news_article_translations").
		Where("slug = ?", slug).
		Where("language = ?", strings.ToUpper(lang))

	if excludeArticleID != nil {
		query = query.Where("article_id <> ?", *excludeArticleID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *newsArticleRepositoryImpl) ReplaceCategories(ctx context.Context, articleID uuid.UUID, categoryIDs []uuid.UUID) error {
	article := &entity.NewsArticle{ID: articleID}

	categories := make([]entity.NewsCategory, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		categories = append(categories, entity.NewsCategory{ID: id})
	}

	return r.db.WithContext(ctx).Model(article).Association("Categories").Replace(&categories)
}

func (r *newsArticleRepositoryImpl) FindCategoriesByIDs(ctx context.Context, categoryIDs []uuid.UUID) ([]entity.NewsCategory, error) {
	if len(categoryIDs) == 0 {
		return []entity.NewsCategory{}, nil
	}
	var categories []entity.NewsCategory
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Where("id IN ?", categoryIDs).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
