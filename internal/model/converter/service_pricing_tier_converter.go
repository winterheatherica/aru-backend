package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func ServicePricingTierToModel(
	tier entity.ServicePricingTier,
	lang string,
) *model.ServicePricingTier {

	tr := findServicePricingTierTranslation(tier, lang)
	if tr == nil {
		return nil
	}

	return &model.ServicePricingTier{
		ID: tier.ID,

		Service: tier.Service,

		PriceMonthly: tier.PriceMonthly,
		PriceYearly:  tier.PriceYearly,

		Popular:    tier.Popular,
		OrderIndex: tier.OrderIndex,

		Name:        tr.Name,
		Description: tr.Description,
		Features:    []string(tr.Features),

		IsActive: tier.IsActive,
	}
}

func findServicePricingTierTranslation(
	tier entity.ServicePricingTier,
	lang string,
) *entity.ServicePricingTierTranslation {

	for i := range tier.Translations {
		if tier.Translations[i].Language == lang {
			return &tier.Translations[i]
		}
	}
	return nil
}

func ServicePricingTierListToModel(
	tiers []entity.ServicePricingTier,
	lang string,
) []model.ServicePricingTier {

	result := make([]model.ServicePricingTier, 0, len(tiers))

	for _, t := range tiers {
		m := ServicePricingTierToModel(t, lang)
		if m != nil {
			result = append(result, *m)
		}
	}

	return result
}
