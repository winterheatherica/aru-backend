package admin

import (
	"context"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ServicePricingTierUsecase interface {
	ListByService(ctx context.Context, service string, lang string) ([]model.ServicePricingTierAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServicePricingTierAdminItem, error)
	Create(ctx context.Context, input model.ServicePricingTierUpsertInput) (*model.ServicePricingTierAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.ServicePricingTierUpsertInput) (*model.ServicePricingTierAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type servicePricingTierUsecaseImpl struct {
	repo repository.ServicePricingRepository
}

func NewServicePricingTierUsecase(repo repository.ServicePricingRepository) ServicePricingTierUsecase {
	return &servicePricingTierUsecaseImpl{repo: repo}
}

func (u *servicePricingTierUsecaseImpl) ListByService(ctx context.Context, service string, lang string) ([]model.ServicePricingTierAdminItem, error) {
	items, err := u.repo.FindByService(ctx, strings.ToUpper(strings.TrimSpace(service)))
	if err != nil {
		return nil, err
	}
	res := make([]model.ServicePricingTierAdminItem, 0, len(items))
	for _, it := range items {
		res = append(res, toServicePricingTierAdminItem(it, lang))
	}
	return res, nil
}

func (u *servicePricingTierUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServicePricingTierAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toServicePricingTierAdminItem(*it, lang)
	return &res, nil
}

func (u *servicePricingTierUsecaseImpl) Create(ctx context.Context, input model.ServicePricingTierUpsertInput) (*model.ServicePricingTierAdminItem, error) {
	lang := normPricingLang(input.Language)
	item := &entity.ServicePricingTier{
		Service:      strings.ToUpper(strings.TrimSpace(input.Service)),
		PriceMonthly: input.PriceMonthly,
		PriceYearly:  input.PriceYearly,
		Popular:      input.Popular,
		OrderIndex:   maxInt(input.OrderIndex, 1),
		IsActive:     input.IsActive,
	}
	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := u.upsertTranslations(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}
	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toServicePricingTierAdminItem(*item, lang)
		return &res, nil
	}
	res := toServicePricingTierAdminItem(*fresh, lang)
	return &res, nil
}

func (u *servicePricingTierUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.ServicePricingTierUpsertInput) (*model.ServicePricingTierAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lang := normPricingLang(input.Language)
	if s := strings.ToUpper(strings.TrimSpace(input.Service)); s != "" {
		item.Service = s
	}
	item.PriceMonthly = input.PriceMonthly
	item.PriceYearly = input.PriceYearly
	item.Popular = input.Popular
	item.OrderIndex = maxInt(input.OrderIndex, 1)
	item.IsActive = input.IsActive
	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	if err := u.upsertTranslations(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}
	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toServicePricingTierAdminItem(*item, lang)
		return &res, nil
	}
	res := toServicePricingTierAdminItem(*fresh, lang)
	return &res, nil
}

func (u *servicePricingTierUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteByID(ctx, id)
}

func (u *servicePricingTierUsecaseImpl) upsertTranslations(ctx context.Context, tierID uuid.UUID, input model.ServicePricingTierUpsertInput, lang string) error {
	if err := u.repo.UpsertTranslation(ctx, &entity.ServicePricingTierTranslation{TierID: tierID, Language: lang, Name: strings.TrimSpace(input.Name), Description: trimPtr(input.Description), Features: toPricingPQStrings(input.Features)}); err != nil {
		return err
	}
	if lang == "ID" {
		enFeatures := make([]string, 0, len(input.Features))
		for _, f := range input.Features {
			enFeatures = append(enFeatures, translateTextValPricing(f, "id", "en"))
		}
		_ = u.repo.UpsertTranslation(ctx, &entity.ServicePricingTierTranslation{TierID: tierID, Language: "EN", Name: translateTextValPricing(input.Name, "id", "en"), Description: translateTextPtrPricing(input.Description, "id", "en"), Features: toPricingPQStrings(enFeatures)})
	}
	return nil
}

func toServicePricingTierAdminItem(it entity.ServicePricingTier, lang string) model.ServicePricingTierAdminItem {
	trs := make([]model.ServicePricingTierAdminTranslation, 0, len(it.Translations))
	for _, tr := range it.Translations {
		trs = append(trs, model.ServicePricingTierAdminTranslation{Language: tr.Language, Name: tr.Name, Description: tr.Description, Features: []string(tr.Features)})
	}
	name := ""
	var desc *string
	features := []string{}
	for _, tr := range it.Translations {
		if strings.EqualFold(tr.Language, lang) {
			name, desc, features = tr.Name, tr.Description, []string(tr.Features)
			break
		}
	}
	if name == "" && len(it.Translations) > 0 {
		name, desc, features = it.Translations[0].Name, it.Translations[0].Description, []string(it.Translations[0].Features)
	}
	return model.ServicePricingTierAdminItem{ID: it.ID, Service: it.Service, PriceMonthly: it.PriceMonthly, PriceYearly: it.PriceYearly, Popular: it.Popular, OrderIndex: it.OrderIndex, IsActive: it.IsActive, Name: name, Description: desc, Features: features, Translations: trs}
}

func normPricingLang(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" {
		return "ID"
	}
	return l
}

func toPricingPQStrings(items []string) pq.StringArray {
	out := make([]string, 0, len(items))
	for _, it := range items {
		if t := strings.TrimSpace(it); t != "" {
			out = append(out, t)
		}
	}
	return pq.StringArray(out)
}

func maxInt(v, fallback int) int {
	if v <= 0 {
		return fallback
	}
	return v
}

func translateTextPtrPricing(text *string, fromLang, toLang string) *string {
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

func translateTextValPricing(text, fromLang, toLang string) string {
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
