package admin

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/minio/minio-go/v7"
)

type NewsArticleUsecase interface {
	List(ctx context.Context) ([]model.NewsArticleAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.NewsArticleAdminItem, error)
	Create(ctx context.Context, input model.NewsArticleUpsertInput, image *multipart.FileHeader) (*model.NewsArticleAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.NewsArticleUpsertInput, image *multipart.FileHeader) (*model.NewsArticleAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type newsArticleUsecaseImpl struct {
	repo          repository.NewsArticleRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewNewsArticleUsecase(repo repository.NewsArticleRepository, minioClient *minio.Client, minioBucket string, publicBaseURL string) NewsArticleUsecase {
	return &newsArticleUsecaseImpl{
		repo:          repo,
		minioClient:   minioClient,
		minioBucket:   minioBucket,
		publicBaseURL: strings.TrimSuffix(publicBaseURL, "/"),
	}
}

func (u *newsArticleUsecaseImpl) List(ctx context.Context) ([]model.NewsArticleAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]model.NewsArticleAdminItem, 0, len(items))
	for _, item := range items {
		adminItem := u.toAdminItem(item)
		u.attachUploaderName(ctx, &adminItem)
		res = append(res, adminItem)
	}
	return res, nil
}

func (u *newsArticleUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.NewsArticleAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*item)
	u.attachUploaderName(ctx, &res)
	return &res, nil
}

func (u *newsArticleUsecaseImpl) Create(ctx context.Context, input model.NewsArticleUpsertInput, image *multipart.FileHeader) (*model.NewsArticleAdminItem, error) {
	lang := normalizeLang(input.Language)
	idTitle := strings.TrimSpace(input.Title)
	if idTitle == "" {
		return nil, fmt.Errorf("title is required")
	}

	idSlug := buildSlugNewsArticle(idTitle)
	idSlug, _ = u.ensureUniqueSlug(ctx, idSlug, lang, nil)

	idContent := strings.TrimSpace(input.Content)
	idMetaTitle := deriveMetaTitleNewsArticle(idTitle)
	idMetaDesc := deriveMetaDescriptionNewsArticle(idContent)

	categoryIDs, parseErr := parseCategoryUUIDs(input.CategoryIDs)
	if parseErr != nil {
		return nil, parseErr
	}
	categories, err := u.repo.FindCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	translations := []entity.NewsArticleTranslation{
		{
			Language:        lang,
			Slug:            idSlug,
			Title:           idTitle,
			Content:         idContent,
			MetaTitle:       idMetaTitle,
			MetaDescription: idMetaDesc,
			MetaKeywords:    toPQKeywords(deriveKeywordsFromCategories(categories, lang)),
		},
	}

	if lang == "ID" {
		enTitle := translateTextNewsArticle(idTitle, "id", "en")
		enSlug := buildSlugNewsArticle(enTitle)
		if enSlug == "" {
			enSlug = buildSlugNewsArticle(idSlug + "-en")
		}
		if enSlug == idSlug {
			enSlug = buildSlugNewsArticle(enSlug + "-en")
		}
		enSlug, _ = u.ensureUniqueSlug(ctx, enSlug, "EN", nil)

		enContent := translateTextNewsArticle(idContent, "id", "en")
		enMetaTitle := deriveMetaTitleNewsArticle(enTitle)
		enMetaDesc := deriveMetaDescriptionNewsArticle(enContent)

		translations = append(translations, entity.NewsArticleTranslation{
			Language:        "EN",
			Slug:            enSlug,
			Title:           enTitle,
			Content:         enContent,
			MetaTitle:       enMetaTitle,
			MetaDescription: enMetaDesc,
			MetaKeywords:    toPQKeywords(deriveKeywordsFromCategories(categories, "EN")),
		})
	}

	article := &entity.NewsArticle{
		IsActive:     input.IsActive,
		PublishedAt:  ensurePublishedAt(input.PublishedAt),
		UploadedBy:   input.UploadedBy,
		Translations: translations,
	}

	if image != nil {
		objectName, err := u.uploadImage(ctx, image)
		if err != nil {
			return nil, err
		}
		article.ImagePath = &objectName
	}

	if err := u.repo.Create(ctx, article); err != nil {
		return nil, err
	}

	if err := u.repo.ReplaceCategories(ctx, article.ID, categoryIDs); err != nil {
		return nil, err
	}

	created, err := u.repo.FindByID(ctx, article.ID)
	if err != nil {
		res := u.toAdminItem(*article)
		u.attachUploaderName(ctx, &res)
		return &res, nil
	}
	res := u.toAdminItem(*created)
	u.attachUploaderName(ctx, &res)
	return &res, nil
}

func (u *newsArticleUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.NewsArticleUpsertInput, image *multipart.FileHeader) (*model.NewsArticleAdminItem, error) {
	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lang := normalizeLang(input.Language)
	slug := buildSlugNewsArticle(input.Title)
	slug, _ = u.ensureUniqueSlug(ctx, slug, lang, &id)

	article.IsActive = input.IsActive
	article.PublishedAt = ensurePublishedAt(input.PublishedAt)
	if input.UploadedBy != nil {
		article.UploadedBy = input.UploadedBy
	}

	oldImagePath := ""
	if article.ImagePath != nil {
		oldImagePath = *article.ImagePath
	}
	if image != nil {
		newObject, upErr := u.uploadImage(ctx, image)
		if upErr != nil {
			return nil, upErr
		}
		article.ImagePath = &newObject
	}

	if err := u.repo.Update(ctx, article); err != nil {
		return nil, err
	}

	trimmedTitle := strings.TrimSpace(input.Title)
	trimmedContent := strings.TrimSpace(input.Content)

	categoryIDs, parseErr := parseCategoryUUIDs(input.CategoryIDs)
	if parseErr != nil {
		return nil, parseErr
	}
	categories, err := u.repo.FindCategoriesByIDs(ctx, categoryIDs)
	if err != nil {
		return nil, err
	}

	tr := &entity.NewsArticleTranslation{
		ArticleID:        article.ID,
		Language:         lang,
		Slug:             slug,
		Title:            trimmedTitle,
		Content:          trimmedContent,
		MetaTitle:        deriveMetaTitleNewsArticle(trimmedTitle),
		MetaDescription:  deriveMetaDescriptionNewsArticle(trimmedContent),
		MetaKeywords:     toPQKeywords(deriveKeywordsFromCategories(categories, lang)),
	}
	if err := u.repo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	if err := u.repo.ReplaceCategories(ctx, article.ID, categoryIDs); err != nil {
		return nil, err
	}

	if refreshed, refErr := u.repo.FindByID(ctx, article.ID); refErr == nil {
		for _, existing := range refreshed.Translations {
			kw := toPQKeywords(deriveKeywordsFromCategories(refreshed.Categories, strings.ToUpper(existing.Language)))
			upsert := &entity.NewsArticleTranslation{
				ArticleID:       article.ID,
				Language:        strings.ToUpper(existing.Language),
				Slug:            existing.Slug,
				Title:           existing.Title,
				Content:         existing.Content,
				MetaTitle:       deriveMetaTitleNewsArticle(existing.Title),
				MetaDescription: deriveMetaDescriptionNewsArticle(existing.Content),
				MetaKeywords:    kw,
			}
			_ = u.repo.UpsertTranslation(ctx, upsert)
		}
	}

	if image != nil && oldImagePath != "" && article.ImagePath != nil && oldImagePath != *article.ImagePath {
		_ = u.removeObject(ctx, oldImagePath)
	}

	updated, err := u.repo.FindByID(ctx, article.ID)
	if err != nil {
		res := u.toAdminItem(*article)
		u.attachUploaderName(ctx, &res)
		return &res, nil
	}
	res := u.toAdminItem(*updated)
	u.attachUploaderName(ctx, &res)
	return &res, nil
}

func (u *newsArticleUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	article, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.repo.DeleteByID(ctx, id); err != nil {
		return err
	}

	if article.ImagePath != nil && strings.TrimSpace(*article.ImagePath) != "" {
		_ = u.removeObject(ctx, *article.ImagePath)
	}
	return nil
}

func (u *newsArticleUsecaseImpl) toAdminItem(item entity.NewsArticle) model.NewsArticleAdminItem {
	translations := make([]model.NewsArticleAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		translations = append(translations, model.NewsArticleAdminTranslation{
			Language:        tr.Language,
			Slug:            tr.Slug,
			Title:           tr.Title,
			Content:         tr.Content,
			MetaTitle:       tr.MetaTitle,
			MetaDescription: tr.MetaDescription,
			MetaKeywords:    []string(tr.MetaKeywords),
		})
	}

	categoryIDs := make([]uuid.UUID, 0, len(item.Categories))
	categories := make([]model.NewsArticleAdminCategory, 0, len(item.Categories))
	for _, c := range item.Categories {
		categoryIDs = append(categoryIDs, c.ID)
		name := ""
		slug := ""
		for _, tr := range c.Translations {
			if strings.EqualFold(tr.Language, "ID") {
				name = tr.Name
				slug = tr.Slug
				break
			}
			if name == "" {
				name = tr.Name
				slug = tr.Slug
			}
		}
		categories = append(categories, model.NewsArticleAdminCategory{ID: c.ID, Name: name, Slug: slug})
	}

	res := model.NewsArticleAdminItem{
		ID:           item.ID,
		IsActive:     item.IsActive,
		PublishedAt:  item.PublishedAt,
		ViewCount:    item.ViewCount,
		LikeCount:    item.LikeCount,
		UploadedBy:   item.UploadedBy,
		CategoryIDs:  categoryIDs,
		Categories:   categories,
		Translations: translations,
	}
	if item.ImagePath != nil {
		res.ImagePath = item.ImagePath
		url := converter.BuildAssetURL(u.publicBaseURL, *item.ImagePath)
		res.ImageURL = &url
	}
	return res
}

func (u *newsArticleUsecaseImpl) ensureUniqueSlug(ctx context.Context, slug, lang string, excludeID *uuid.UUID) (string, error) {
	base := strings.TrimSpace(slug)
	if base == "" {
		base = "news-article"
	}

	candidate := base
	for i := 1; i <= 200; i++ {
		exists, err := u.repo.IsSlugExists(ctx, candidate, lang, excludeID)
		if err != nil {
			return candidate, err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}

	return fmt.Sprintf("%s-%d", base, time.Now().Unix()), nil
}

func ensurePublishedAt(t time.Time) time.Time {
	if t.IsZero() {
		return time.Now().UTC()
	}
	return t.UTC()
}

func normalizeLang(lang string) string {
	res := strings.ToUpper(strings.TrimSpace(lang))
	if res == "" {
		return "ID"
	}
	return res
}

var slugNonAlnumNewsArticle = regexp.MustCompile(`[^a-z0-9]+`)

func buildSlugNewsArticle(text string) string {
	t := strings.ToLower(strings.TrimSpace(text))
	t = slugNonAlnumNewsArticle.ReplaceAllString(t, "-")
	t = strings.Trim(t, "-")
	if t == "" {
		return "news-article"
	}
	return t
}

func toPQKeywords(values []string) pq.StringArray {
	if len(values) == 0 {
		return nil
	}
	res := make([]string, 0, len(values))
	for _, v := range values {
		t := strings.TrimSpace(v)
		if t != "" {
			res = append(res, t)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return pq.StringArray(res)
}

func deriveKeywordsFromCategories(categories []entity.NewsCategory, lang string) []string {
	if len(categories) == 0 {
		return []string{}
	}
	upLang := strings.ToUpper(strings.TrimSpace(lang))
	seen := map[string]bool{}
	result := make([]string, 0, len(categories))

	for _, c := range categories {
		picked := ""
		for _, tr := range c.Translations {
			if strings.ToUpper(strings.TrimSpace(tr.Language)) == upLang {
				picked = strings.TrimSpace(tr.Name)
				break
			}
		}
		if picked == "" {
			for _, tr := range c.Translations {
				if strings.EqualFold(tr.Language, "ID") {
					picked = strings.TrimSpace(tr.Name)
					break
				}
			}
		}
		if picked == "" && len(c.Translations) > 0 {
			picked = strings.TrimSpace(c.Translations[0].Name)
		}
		if picked != "" && !seen[picked] {
			seen[picked] = true
			result = append(result, picked)
		}
	}

	return result
}

func (u *newsArticleUsecaseImpl) attachUploaderName(ctx context.Context, item *model.NewsArticleAdminItem) {
	if item == nil || item.UploadedBy == nil {
		return
	}
	publisher, err := u.repo.FindPublisherByArticleID(ctx, item.ID.String())
	if err != nil || publisher == nil {
		return
	}
	name := ""
	if publisher.FullName != nil {
		name = strings.TrimSpace(*publisher.FullName)
	}
	if name == "" && publisher.Username != nil {
		name = strings.TrimSpace(*publisher.Username)
	}
	if name == "" {
		name = strings.TrimSpace(publisher.Email)
	}
	if name == "" {
		return
	}
	item.UploadedByName = &name
}

func translateTextNewsArticle(text, fromLang, toLang string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}
	translated, err := gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{From: fromLang, To: toLang})
	if err != nil || strings.TrimSpace(translated) == "" {
		return trimmed
	}
	return strings.TrimSpace(translated)
}

func deriveMetaTitleNewsArticle(title string) *string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return nil
	}
	runes := []rune(trimmed)
	if len(runes) > 70 {
		trimmed = strings.TrimSpace(string(runes[:70]))
	}
	return &trimmed
}

var htmlTagRegexNewsArticle = regexp.MustCompile(`<[^>]+>`)

func deriveMetaDescriptionNewsArticle(content string) *string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return nil
	}
	plain := htmlTagRegexNewsArticle.ReplaceAllString(trimmed, " ")
	plain = strings.Join(strings.Fields(plain), " ")
	if plain == "" {
		return nil
	}
	runes := []rune(plain)
	if len(runes) > 160 {
		plain = strings.TrimSpace(string(runes[:160]))
	}
	return &plain
}

func translateKeywordsNewsArticle(values []string, fromLang, toLang string) pq.StringArray {
	if len(values) == 0 {
		return nil
	}
	translated := make([]string, 0, len(values))
	for _, value := range values {
		v := translateTextNewsArticle(value, fromLang, toLang)
		if strings.TrimSpace(v) != "" {
			translated = append(translated, v)
		}
	}
	if len(translated) == 0 {
		return nil
	}
	return pq.StringArray(translated)
}

func parseCategoryUUIDs(values []string) ([]uuid.UUID, error) {
	if len(values) == 0 {
		return []uuid.UUID{}, nil
	}
	res := make([]uuid.UUID, 0, len(values))
	seen := map[uuid.UUID]bool{}
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		id, err := uuid.Parse(trimmed)
		if err != nil {
			return nil, fmt.Errorf("invalid category_id: %s", trimmed)
		}
		if !seen[id] {
			seen[id] = true
			res = append(res, id)
		}
	}
	return res, nil
}

func (u *newsArticleUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/news_articles/%s%s", uuid.NewString(), ext)

	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{
		ContentType: image.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *newsArticleUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
