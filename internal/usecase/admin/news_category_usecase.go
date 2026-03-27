package admin

import (
	"context"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
)

type NewsCategoryUsecase interface {
	List(ctx context.Context) ([]model.NewsCategoryAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.NewsCategoryAdminItem, error)
	Create(ctx context.Context, input model.NewsCategoryUpsertInput) (*model.NewsCategoryAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.NewsCategoryUpsertInput) (*model.NewsCategoryAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type newsCategoryUsecaseImpl struct {
	repo repository.NewsCategoryRepository
}

func NewNewsCategoryUsecase(repo repository.NewsCategoryRepository) NewsCategoryUsecase {
	return &newsCategoryUsecaseImpl{repo: repo}
}

func (u *newsCategoryUsecaseImpl) List(ctx context.Context) ([]model.NewsCategoryAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.NewsCategoryAdminItem, 0, len(items))
	for _, item := range items {
		result = append(result, toNewsCategoryAdminItem(item))
	}

	return result, nil
}

func (u *newsCategoryUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.NewsCategoryAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := toNewsCategoryAdminItem(*item)
	return &res, nil
}

func (u *newsCategoryUsecaseImpl) Create(ctx context.Context, input model.NewsCategoryUpsertInput) (*model.NewsCategoryAdminItem, error) {
	language := strings.ToUpper(strings.TrimSpace(input.Language))
	if language == "" {
		language = "ID"
	}

	name := strings.TrimSpace(input.Name)
	slug := strings.TrimSpace(input.Slug)

	translations := []entity.NewsCategoryTranslation{
		{
			Language:    language,
			Name:        name,
			Slug:        slug,
			Description: input.Description,
		},
	}

	if language == "ID" {
		translations = append(translations, entity.NewsCategoryTranslation{
			Language:    "EN",
			Name:        translateTextNewsCategory(name, "id", "en"),
			Slug:        slug,
			Description: translateTextPtrNewsCategory(input.Description, "id", "en"),
		})
	}

	item := &entity.NewsCategory{
		IsActive:     input.IsActive,
		Translations: translations,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	created, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toNewsCategoryAdminItem(*item)
		return &res, nil
	}

	res := toNewsCategoryAdminItem(*created)
	return &res, nil
}

func (u *newsCategoryUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.NewsCategoryUpsertInput) (*model.NewsCategoryAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	language := strings.ToUpper(strings.TrimSpace(input.Language))
	if language == "" {
		language = "ID"
	}

	item.IsActive = input.IsActive
	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.NewsCategoryTranslation{
		CategoryID:  item.ID,
		Language:    language,
		Name:        strings.TrimSpace(input.Name),
		Slug:        strings.TrimSpace(input.Slug),
		Description: input.Description,
	}
	if err := u.repo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	updated, err := u.repo.FindByID(ctx, id)
	if err != nil {
		res := toNewsCategoryAdminItem(*item)
		return &res, nil
	}

	res := toNewsCategoryAdminItem(*updated)
	return &res, nil
}

func (u *newsCategoryUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteByID(ctx, id)
}

func toNewsCategoryAdminItem(item entity.NewsCategory) model.NewsCategoryAdminItem {
	translations := make([]model.NewsCategoryAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		translations = append(translations, model.NewsCategoryAdminTranslation{
			Language:    tr.Language,
			Name:        tr.Name,
			Slug:        tr.Slug,
			Description: tr.Description,
		})
	}

	res := model.NewsCategoryAdminItem{
		ID:           item.ID,
		IsActive:     item.IsActive,
		Translations: translations,
	}

	if len(item.Translations) > 0 {
		first := item.Translations[0]
		res.Name = &first.Name
		res.Slug = &first.Slug
		res.Description = first.Description
	}

	return res
}

func translateTextNewsCategory(text string, fromLang, toLang string) string {
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

func translateTextPtrNewsCategory(text *string, fromLang, toLang string) *string {
	if text == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*text)
	if trimmed == "" {
		return nil
	}

	translated, err := gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{From: fromLang, To: toLang})
	if err != nil || strings.TrimSpace(translated) == "" {
		fallback := trimmed
		return &fallback
	}

	res := strings.TrimSpace(translated)
	return &res
}
