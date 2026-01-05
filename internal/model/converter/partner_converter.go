package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func PartnerToModel(
	partner entity.Partner,
	lang string,
	baseURL string,
) *model.Partner {

	trans := findPartnerTranslation(partner, lang)

	alt := ""
	var title *string
	var desc *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
		desc = trans.Description
	}

	return &model.Partner{
		ID:    partner.ID,
		Src:   BuildAssetURL(baseURL, partner.ImagePath),
		Alt:   alt,
		Title: title,
		Desc:  desc,
		Order: partner.OrderIndex,
	}
}

func PartnerListToModel(
	partners []entity.Partner,
	lang string,
	baseURL string,
) []model.Partner {

	result := make([]model.Partner, 0, len(partners))

	for _, p := range partners {
		result = append(result, *PartnerToModel(p, lang, baseURL))
	}

	return result
}

func findPartnerTranslation(
	partner entity.Partner,
	lang string,
) *entity.PartnerTranslation {

	for i := range partner.Translations {
		if partner.Translations[i].Language == lang {
			return &partner.Translations[i]
		}
	}
	return nil
}
