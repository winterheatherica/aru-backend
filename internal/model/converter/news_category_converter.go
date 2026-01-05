package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func NewsCategoryToModel(
	category entity.NewsCategory,
	lang string,
) *model.NewsCategory {

	trans := findNewsCategoryTranslation(category, lang)

	if trans == nil {
		return nil
	}

	return &model.NewsCategory{
		ID:          category.ID,
		Name:        trans.Name,
		Slug:        trans.Slug,
		Description: trans.Description,
		Meta: model.NewsCategoryMeta{
			Title:       trans.MetaTitle,
			Description: trans.MetaDescription,
			Keywords:    trans.MetaKeywords,
		},
	}
}

func findNewsCategoryTranslation(
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
