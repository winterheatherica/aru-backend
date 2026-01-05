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
	var trans *entity.PromoSlideTranslation

	for i := range slide.Translations {
		if slide.Translations[i].Language == lang {
			trans = &slide.Translations[i]
			break
		}
	}

	alt := ""
	var title *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
	}

	src := buildMinioURL(baseURL, slide.ImagePath)

	return &model.PromoSlide{
		ID:    slide.ID,
		Src:   src,
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

	for _, s := range slides {
		m := PromoSlideToModel(s, lang, baseURL)
		if m != nil {
			result = append(result, *m)
		}
	}

	return result
}
