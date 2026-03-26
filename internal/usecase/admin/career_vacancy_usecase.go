package admin

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
)

type CareerVacancyUsecase interface {
	List(ctx context.Context) ([]model.CareerVacancyAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.CareerVacancyAdminItem, error)
	Create(ctx context.Context, input model.CareerVacancyUpsertInput) (*model.CareerVacancyAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.CareerVacancyUpsertInput) (*model.CareerVacancyAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type careerVacancyUsecaseImpl struct {
	repo repository.CareerVacancyRepository
}

func NewCareerVacancyUsecase(repo repository.CareerVacancyRepository) CareerVacancyUsecase {
	return &careerVacancyUsecaseImpl{repo: repo}
}

func (u *careerVacancyUsecaseImpl) List(ctx context.Context) ([]model.CareerVacancyAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.CareerVacancyAdminItem, 0, len(items))
	for _, it := range items {
		result = append(result, toCareerAdminItem(it))
	}
	return result, nil
}

func (u *careerVacancyUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.CareerVacancyAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toCareerAdminItem(*it)
	return &res, nil
}

func (u *careerVacancyUsecaseImpl) Create(ctx context.Context, input model.CareerVacancyUpsertInput) (*model.CareerVacancyAdminItem, error) {
	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	openedAt, closedAt, err := parseOpenClose(input.OpenedAt, input.ClosedAt)
	if err != nil {
		return nil, err
	}

	item := &entity.CareerVacancy{
		Title:      title,
		Slug:       buildSlug(title),
		Employment: strings.ToUpper(strings.TrimSpace(input.Employment)),
		Location:   strings.TrimSpace(input.Location),
		IsActive:   input.IsActive,
		OpenedAt:   openedAt,
		ClosedAt:   closedAt,
		Translations: []entity.CareerVacancyTranslation{
			{
				Language:    strings.ToUpper(input.Language),
				Description: input.Description,
			},
		},
	}

	if strings.EqualFold(input.Language, "ID") {
		item.Translations = append(item.Translations, entity.CareerVacancyTranslation{
			Language:    "EN",
			Description: translateTextPtrCareer(input.Description, "id", "en"),
		})
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	created, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toCareerAdminItem(*item)
		return &res, nil
	}
	res := toCareerAdminItem(*created)
	return &res, nil
}

func (u *careerVacancyUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.CareerVacancyUpsertInput) (*model.CareerVacancyAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	openedAt, closedAt, err := parseOpenClose(input.OpenedAt, input.ClosedAt)
	if err != nil {
		return nil, err
	}

	item.Title = title
	item.Slug = buildSlug(title)
	item.Employment = strings.ToUpper(strings.TrimSpace(input.Employment))
	item.Location = strings.TrimSpace(input.Location)
	item.IsActive = input.IsActive
	item.OpenedAt = openedAt
	item.ClosedAt = closedAt

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.CareerVacancyTranslation{
		VacancyID:   item.ID,
		Language:    strings.ToUpper(input.Language),
		Description: input.Description,
	}
	if err := u.repo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	updated, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toCareerAdminItem(*item)
		return &res, nil
	}
	res := toCareerAdminItem(*updated)
	return &res, nil
}

func (u *careerVacancyUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	_, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return u.repo.DeleteByID(ctx, id)
}

func toCareerAdminItem(item entity.CareerVacancy) model.CareerVacancyAdminItem {
	translations := make([]model.CareerVacancyAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		translations = append(translations, model.CareerVacancyAdminTranslation{
			Language:    tr.Language,
			Description: tr.Description,
		})
	}

	opened := item.OpenedAt.Format(time.RFC3339)
	var closed *string
	if item.ClosedAt != nil {
		v := item.ClosedAt.Format(time.RFC3339)
		closed = &v
	}

	return model.CareerVacancyAdminItem{
		ID:           item.ID,
		Title:        item.Title,
		Slug:         item.Slug,
		Employment:   item.Employment,
		Location:     item.Location,
		IsActive:     item.IsActive,
		OpenedAt:     opened,
		ClosedAt:     closed,
		Translations: translations,
	}
}

func parseOpenClose(openedRaw string, closedRaw *string) (time.Time, *time.Time, error) {
	openedRaw = strings.TrimSpace(openedRaw)
	if openedRaw == "" {
		return time.Time{}, nil, fmt.Errorf("opened_at is required")
	}

	openedAt, err := time.Parse("2006-01-02", openedRaw)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("opened_at must be YYYY-MM-DD")
	}

	if closedRaw == nil || strings.TrimSpace(*closedRaw) == "" {
		return openedAt, nil, nil
	}
	closedAt, cerr := time.Parse("2006-01-02", strings.TrimSpace(*closedRaw))
	if cerr != nil {
		return time.Time{}, nil, fmt.Errorf("closed_at must be YYYY-MM-DD")
	}
	return openedAt, &closedAt, nil
}

var slugNonAlnum = regexp.MustCompile(`[^a-z0-9]+`)

func buildSlug(title string) string {
	t := strings.ToLower(strings.TrimSpace(title))
	t = slugNonAlnum.ReplaceAllString(t, "-")
	t = strings.Trim(t, "-")
	if t == "" {
		return uuid.NewString()
	}
	return t
}

func translateTextPtrCareer(text *string, fromLang, toLang string) *string {
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
