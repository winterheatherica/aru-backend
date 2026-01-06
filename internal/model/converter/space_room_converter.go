package converter

import (
	"time"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func SpaceRoomListToCardModel(
	rooms []entity.SpaceRoom,
	lang string,
	baseURL string,
) []model.SpaceRoomCard {

	result := make([]model.SpaceRoomCard, 0, len(rooms))

	for _, r := range rooms {
		card := SpaceRoomToCardModel(r, lang, baseURL)
		if card != nil {
			result = append(result, *card)
		}
	}

	return result
}

func SpaceRoomToCardModel(
	room entity.SpaceRoom,
	lang string,
	baseURL string,
) *model.SpaceRoomCard {

	tr := findSpaceRoomTranslation(room, lang)
	if tr == nil {
		return nil
	}

	var mainImageURL *string
	if room.MainImageURL != nil {
		url := BuildAssetURL(baseURL, *room.MainImageURL)
		mainImageURL = &url
	}

	isAvailable := isRoomAvailable(room.Bookings, time.Now())
	statusText := "Available"
	actionLabel := "View & Book"
	actionState := "active"

	if !isAvailable {
		statusText = "Fully booked"
		actionLabel = "Unavailable"
		actionState = "disabled"
	}

	return &model.SpaceRoomCard{
		ID: room.ID,

		Slug: tr.Slug,

		Title:       tr.Title,
		Description: tr.Description,

		MainImageURL:   mainImageURL,
		MainImageAlt:   tr.MainImageAlt,
		MainImageTitle: tr.MainImageTitle,

		Capacity: room.Capacity,
		Floor:    room.Floor,

		Facilities: tr.Facilities,

		IsAvailable: isAvailable,
		StatusText:  statusText,

		Price: nil,

		ActionLabel: actionLabel,
		ActionState: actionState,
	}
}

func findSpaceRoomTranslation(
	room entity.SpaceRoom,
	lang string,
) *entity.SpaceRoomTranslation {

	for i := range room.Translations {
		if room.Translations[i].Language == lang {
			return &room.Translations[i]
		}
	}
	return nil
}

func isRoomAvailable(
	bookings []entity.SpaceRoomBooking,
	now time.Time,
) bool {

	for _, b := range bookings {
		if b.Status != "CANCELLED" &&
			now.After(b.StartTime) &&
			now.Before(b.EndTime) {
			return false
		}
	}
	return true
}

func SpaceRoomToDetailModel(
	room entity.SpaceRoom,
	lang string,
	baseURL string,
) *model.SpaceRoomDetail {

	tr := findSpaceRoomTranslation(room, lang)
	if tr == nil {
		return nil
	}

	// --- images ---
	images := make([]model.SpaceRoomImage, 0)

	for _, img := range room.Images {
		it := findSpaceRoomImageTranslation(img, lang)

		url := BuildAssetURL(baseURL, img.ImageURL)

		var alt *string
		var title *string

		if it != nil {
			alt = it.Alt
			title = it.Title
		}

		images = append(images, model.SpaceRoomImage{
			URL:   url,
			Alt:   alt,
			Title: title,
		})
	}

	isAvailable := isRoomAvailable(room.Bookings, time.Now())

	statusText := "Available"
	actionLabel := "View & Book"
	actionState := "active"

	if !isAvailable {
		statusText = "Fully booked"
		actionLabel = "Unavailable"
		actionState = "disabled"
	}

	return &model.SpaceRoomDetail{
		ID: room.ID,

		Slug: tr.Slug,

		Title:       tr.Title,
		Description: tr.Description,

		Images: images,

		Capacity: room.Capacity,
		Floor:    room.Floor,

		Facilities: tr.Facilities,

		IsAvailable: isAvailable,
		StatusText:  statusText,
		Price:       nil,

		ActionLabel: actionLabel,
		ActionState: actionState,
	}
}

func findSpaceRoomImageTranslation(
	image entity.SpaceRoomImage,
	lang string,
) *entity.SpaceRoomImageTranslation {

	for i := range image.Translations {
		if image.Translations[i].Language == lang {
			return &image.Translations[i]
		}
	}
	return nil
}
