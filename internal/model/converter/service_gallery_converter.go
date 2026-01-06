package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func ServiceGalleryToModel(
	gallery entity.ServiceGallery,
	lang string,
	baseURL string,
) *model.ServiceGallery {

	tr := findServiceGalleryTranslation(gallery, lang)
	if tr == nil {
		return nil
	}

	src := BuildAssetURL(baseURL, gallery.Src)

	var thumbnail *string
	if gallery.Thumbnail != nil {
		t := BuildAssetURL(baseURL, *gallery.Thumbnail)
		thumbnail = &t
	}

	return &model.ServiceGallery{
		ID: gallery.ID,

		Service:   gallery.Service,
		MediaType: gallery.MediaType,

		Src:       src,
		Thumbnail: thumbnail,

		Title:   tr.Title,
		Alt:     tr.Alt,
		Caption: tr.Caption,
	}
}

func findServiceGalleryTranslation(
	gallery entity.ServiceGallery,
	lang string,
) *entity.ServiceGalleryTranslation {

	for i := range gallery.Translations {
		if gallery.Translations[i].Language == lang {
			return &gallery.Translations[i]
		}
	}
	return nil
}

func ServiceGalleryListToModel(
	galleries []entity.ServiceGallery,
	lang string,
	baseURL string,
) []model.ServiceGallery {

	result := make([]model.ServiceGallery, 0, len(galleries))

	for _, g := range galleries {
		m := ServiceGalleryToModel(g, lang, baseURL)
		if m != nil {
			result = append(result, *m)
		}
	}

	return result
}
