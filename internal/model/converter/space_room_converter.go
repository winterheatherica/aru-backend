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

		Rating:     nil,
		RatingText: nil,

		IsAvailable: isAvailable,
		StatusText:  statusText,

		Price: nil,

		Tags: []string{},

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
