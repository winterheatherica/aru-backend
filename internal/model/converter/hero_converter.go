package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"strings"
)

func buildMinioURL(baseURL, objectKey string) string {
	key := strings.TrimPrefix(objectKey, "/")
	return strings.TrimRight(baseURL, "/") + "/" + key
}

func HeroSlideToModel(slide entity.HeroSlide, lang string, baseURL string) *model.HeroSlide {
	var trans *entity.HeroSlideTranslation

	for i := range slide.Translations {
		if slide.Translations[i].Language == lang {
			trans = &slide.Translations[i]
			break
		}
	}

	alt := ""
	var title *string
	var ctaLabel *string
	banner := "POLISH"

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
		ctaLabel = trans.CtaLabel
		banner = trans.Banner
	}

	src := buildMinioURL(baseURL, slide.ImagePath)

	return &model.HeroSlide{
		ID:       slide.ID,
		Src:      src,
		Alt:      alt,
		Title:    title,
		CtaLabel: ctaLabel,
		Banner:   banner,
		Order:    slide.OrderIndex,
	}
}

func HeroSlideListToModel(slides []entity.HeroSlide, lang string, baseURL string) []model.HeroSlide {
	result := make([]model.HeroSlide, 0, len(slides))
	for _, s := range slides {
		m := HeroSlideToModel(s, lang, baseURL)
		if m != nil {
			result = append(result, *m)
		}
	}
	return result
}
