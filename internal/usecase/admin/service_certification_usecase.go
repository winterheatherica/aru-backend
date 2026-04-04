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

type ServiceCertificationUsecase interface {
	ListByService(ctx context.Context, service string, lang string) ([]model.ServiceCertificationAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceCertificationAdminItem, error)
	Create(ctx context.Context, input model.ServiceCertificationUpsertInput) (*model.ServiceCertificationAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.ServiceCertificationUpsertInput) (*model.ServiceCertificationAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type serviceCertificationUsecaseImpl struct {
	repo repository.ServiceCertificationRepository
}

func NewServiceCertificationUsecase(repo repository.ServiceCertificationRepository) ServiceCertificationUsecase {
	return &serviceCertificationUsecaseImpl{repo: repo}
}

func (u *serviceCertificationUsecaseImpl) ListByService(ctx context.Context, service string, lang string) ([]model.ServiceCertificationAdminItem, error) {
	items, err := u.repo.FindByService(ctx, strings.ToUpper(strings.TrimSpace(service)))
	if err != nil {
		return nil, err
	}
	res := make([]model.ServiceCertificationAdminItem, 0, len(items))
	for _, it := range items {
		res = append(res, toServiceCertAdminItem(it, lang))
	}
	return res, nil
}

func (u *serviceCertificationUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceCertificationAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toServiceCertAdminItem(*it, lang)
	return &res, nil
}

func (u *serviceCertificationUsecaseImpl) Create(ctx context.Context, input model.ServiceCertificationUpsertInput) (*model.ServiceCertificationAdminItem, error) {
	service := strings.ToUpper(strings.TrimSpace(input.Service))
	lang := normServiceCertLang(input.Language)

	item := &entity.ServiceCertification{
		Service:    service,
		OrderIndex: input.OrderIndex,
		IsActive:   input.IsActive,
	}
	if item.OrderIndex <= 0 {
		item.OrderIndex = 1
	}
	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if err := u.repo.UpsertTranslation(ctx, &entity.ServiceCertificationTranslation{
		CertificationID: item.ID,
		Language:        lang,
		Title:           trimPtr(input.Title),
		Alt:             trimPtr(input.Alt),
		Caption:         trimPtr(input.Caption),
	}); err != nil {
		return nil, err
	}

	if lang == "ID" {
		enTitle := translateTextPtrServiceCert(input.Title, "id", "en")
		enAlt := translateTextPtrServiceCert(input.Alt, "id", "en")
		enCaption := translateTextPtrServiceCert(input.Caption, "id", "en")
		_ = u.repo.UpsertTranslation(ctx, &entity.ServiceCertificationTranslation{
			CertificationID: item.ID,
			Language:        "EN",
			Title:           enTitle,
			Alt:             enAlt,
			Caption:         enCaption,
		})
	}

	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toServiceCertAdminItem(*item, lang)
		return &res, nil
	}
	res := toServiceCertAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceCertificationUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.ServiceCertificationUpsertInput) (*model.ServiceCertificationAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lang := normServiceCertLang(input.Language)
	if s := strings.ToUpper(strings.TrimSpace(input.Service)); s != "" {
		item.Service = s
	}
	if input.OrderIndex > 0 {
		item.OrderIndex = input.OrderIndex
	}
	item.IsActive = input.IsActive
	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	if err := u.repo.UpsertTranslation(ctx, &entity.ServiceCertificationTranslation{
		CertificationID: item.ID,
		Language:        lang,
		Title:           trimPtr(input.Title),
		Alt:             trimPtr(input.Alt),
		Caption:         trimPtr(input.Caption),
	}); err != nil {
		return nil, err
	}

	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toServiceCertAdminItem(*item, lang)
		return &res, nil
	}
	res := toServiceCertAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceCertificationUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteByID(ctx, id)
}

func normServiceCertLang(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" {
		return "ID"
	}
	return l
}

func toServiceCertAdminItem(it entity.ServiceCertification, lang string) model.ServiceCertificationAdminItem {
	translations := make([]model.ServiceCertificationAdminTranslation, 0, len(it.Translations))
	for _, tr := range it.Translations {
		translations = append(translations, model.ServiceCertificationAdminTranslation{
			Language: tr.Language,
			Title:    tr.Title,
			Alt:      tr.Alt,
			Caption:  tr.Caption,
		})
	}
	var title, alt, caption *string
	for _, tr := range it.Translations {
		if strings.EqualFold(tr.Language, lang) {
			title, alt, caption = tr.Title, tr.Alt, tr.Caption
			break
		}
	}
	if title == nil && len(it.Translations) > 0 {
		title, alt, caption = it.Translations[0].Title, it.Translations[0].Alt, it.Translations[0].Caption
	}
	return model.ServiceCertificationAdminItem{
		ID:           it.ID,
		Service:      it.Service,
		OrderIndex:   it.OrderIndex,
		IsActive:     it.IsActive,
		Title:        title,
		Alt:          alt,
		Caption:      caption,
		Translations: translations,
	}
}

func trimPtr(v *string) *string {
	if v == nil {
		return nil
	}
	t := strings.TrimSpace(*v)
	if t == "" {
		return nil
	}
	return &t
}

func translateTextPtrServiceCert(text *string, fromLang, toLang string) *string {
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
