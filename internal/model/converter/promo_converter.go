package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func PromoSlideToModel(
	slide entity.PromoSlide,
	lang string,
	baseURL string,
) *model.PromoSlide {

	trans := findPromoTranslation(slide, lang)

	alt := ""
	var title *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
	}

	return &model.PromoSlide{
		ID:    slide.ID,
		Src:   BuildAssetURL(baseURL, slide.ImagePath),
		Alt:   alt,
		Title: title,
		Order: slide.OrderIndex,
	}
}

func PromoSlideListToModel(
	slides []entity.PromoSlide,
	lang string,
	baseURL string,
) []model.PromoSlide {

	result := make([]model.PromoSlide, 0, len(slides))

	for _, slide := range slides {
		result = append(result, *PromoSlideToModel(slide, lang, baseURL))
	}

	return result
}

func findPromoTranslation(
	slide entity.PromoSlide,
	lang string,
) *entity.PromoSlideTranslation {

	for i := range slide.Translations {
		if slide.Translations[i].Language == lang {
			return &slide.Translations[i]
		}
	}
	return nil
}
