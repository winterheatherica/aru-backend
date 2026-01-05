package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func HeroSlideToModel(
	slide entity.HeroSlide,
	lang string,
	baseURL string,
) *model.HeroSlide {

	trans := findHeroTranslation(slide, lang)

	alt := ""
	var title *string
	var ctaLabel *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
		ctaLabel = trans.CtaLabel
	}

	return &model.HeroSlide{
		ID:       slide.ID,
		Src:      BuildAssetURL(baseURL, slide.ImagePath),
		Alt:      alt,
		Title:    title,
		CtaLabel: ctaLabel,
		Banner:   slide.Banner,
		Order:    slide.OrderIndex,
	}
}

func HeroSlideListToModel(
	slides []entity.HeroSlide,
	lang string,
	baseURL string,
) []model.HeroSlide {

	result := make([]model.HeroSlide, 0, len(slides))

	for _, slide := range slides {
		result = append(result, *HeroSlideToModel(slide, lang, baseURL))
	}

	return result
}

func findHeroTranslation(
	slide entity.HeroSlide,
	lang string,
) *entity.HeroSlideTranslation {

	for i := range slide.Translations {
		if slide.Translations[i].Language == lang {
			return &slide.Translations[i]
		}
	}
	return nil
}
