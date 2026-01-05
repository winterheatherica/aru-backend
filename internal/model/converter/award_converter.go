package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func AwardToModel(
	award entity.Award,
	lang string,
	baseURL string,
) *model.Award {

	trans := findAwardTranslation(award, lang)

	alt := ""
	var title *string
	var label *string
	var desc *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
		label = trans.Label
		desc = trans.Description
	}

	return &model.Award{
		ID:          award.ID,
		Src:         BuildAssetURL(baseURL, award.ImagePath),
		Alt:         alt,
		Title:       title,
		Label:       label,
		Description: desc,
		Year:        award.Year,
		Order:       award.OrderIndex,
	}
}

func AwardListToModel(
	awards []entity.Award,
	lang string,
	baseURL string,
) []model.Award {

	result := make([]model.Award, 0, len(awards))

	for _, a := range awards {
		result = append(result, *AwardToModel(a, lang, baseURL))
	}

	return result
}

func findAwardTranslation(
	award entity.Award,
	lang string,
) *entity.AwardTranslation {

	for i := range award.Translations {
		if award.Translations[i].Language == lang {
			return &award.Translations[i]
		}
	}
	return nil
}
