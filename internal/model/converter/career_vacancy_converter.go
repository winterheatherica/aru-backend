package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func CareerVacancyToModel(
	v entity.CareerVacancy,
	lang string,
) *model.CareerVacancy {

	desc := ""

	for i := range v.Translations {
		if v.Translations[i].Language == lang {
			if v.Translations[i].Description != nil {
				desc = *v.Translations[i].Description
			}
			break
		}
	}

	return &model.CareerVacancy{
		ID:          v.ID,
		Title:       v.Title,
		Employment:  v.Employment,
		Location:    v.Location,
		Description: desc,
	}
}

func CareerVacancyListToModel(
	vacancies []entity.CareerVacancy,
	lang string,
) []model.CareerVacancy {

	result := make([]model.CareerVacancy, 0, len(vacancies))

	for _, v := range vacancies {
		result = append(result, *CareerVacancyToModel(v, lang))
	}

	return result
}
