package converter

import (
	"regexp"
	"strings"
	"time"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func NewsArticleToModel(
	article entity.NewsArticle,
	lang string,
	baseURL string,
) *model.NewsArticle {

	at := findArticleTranslation(article, lang)
	if at == nil {
		return nil
	}

	var imageURL *string
	if article.ImagePath != nil {
		url := BuildAssetURL(baseURL, *article.ImagePath)
		imageURL = &url
	}

	categories := make([]model.NewsArticleCategory, 0)

	for _, c := range article.Categories {
		ct := findCategoryTranslation(c, lang)
		if ct == nil {
			continue
		}

		categories = append(categories, model.NewsArticleCategory{
			ID:          c.ID,
			Name:        ct.Name,
			Slug:        ct.Slug,
			Description: ct.Description,
		})
	}

	return &model.NewsArticle{
		ID:       article.ID,
		ImageURL: imageURL,

		IsActive:    article.IsActive,
		PublishedAt: article.PublishedAt.Format(time.RFC3339),

		ViewCount: article.ViewCount,
		LikeCount: article.LikeCount,

		CreatedAt: article.CreatedAt.Format(time.RFC3339),
		UpdatedAt: article.UpdatedAt.Format(time.RFC3339),

		Slug:    at.Slug,
		Title:   at.Title,
		Content: at.Content,

		Meta: model.NewsArticleMeta{
			Title:       at.MetaTitle,
			Description: at.MetaDescription,
			Keywords:    at.MetaKeywords,
		},

		Categories: categories,
	}
}

func NewsArticleToNewsCard(
	article entity.NewsArticle,
	lang string,
	baseURL string,
) *model.NewsCard {

	at := findArticleTranslation(article, lang)
	if at == nil {
		return nil
	}

	var imageURL *string
	if article.ImagePath != nil {
		url := BuildAssetURL(baseURL, *article.ImagePath)
		imageURL = &url
	}

	summary := makeSummary(at.Content, 150)

	return &model.NewsCard{
		ID:          article.ID,
		Slug:        at.Slug,
		Title:       at.Title,
		Summary:     summary,
		ImagePath:   imageURL,
		PublishedAt: article.PublishedAt.Format(time.RFC3339),
	}
}

func makeSummary(content string, limit int) string {
	plain := stripHTML(content)
	text := strings.TrimSpace(plain)

	if len(text) <= limit {
		return text
	}
	return text[:limit] + "..."
}

func findArticleTranslation(
	article entity.NewsArticle,
	lang string,
) *entity.NewsArticleTranslation {
	for i := range article.Translations {
		if article.Translations[i].Language == lang {
			return &article.Translations[i]
		}
	}
	return nil
}

func findCategoryTranslation(
	category entity.NewsCategory,
	lang string,
) *entity.NewsCategoryTranslation {
	for i := range category.Translations {
		if category.Translations[i].Language == lang {
			return &category.Translations[i]
		}
	}
	return nil
}

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

func stripHTML(input string) string {
	return htmlTagRe.ReplaceAllString(input, "")
}
