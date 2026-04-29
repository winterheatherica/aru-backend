package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type ReservationUsecase interface {
	GetReservationPage(ctx context.Context, lang string) (*ReservationResponse, error)
}

type reservationUsecaseImpl struct {
	roomRepo repository.SpaceRoomRepository
	baseURL  string
}

type ReservationResponse struct {
	Rooms []model.SpaceRoomCard `json:"rooms"`
}

func NewReservationUsecase(
	roomRepo repository.SpaceRoomRepository,
	baseURL string,
) ReservationUsecase {
	return &reservationUsecaseImpl{
		roomRepo: roomRepo,
		baseURL:  baseURL,
	}
}

func (u *reservationUsecaseImpl) GetReservationPage(
	ctx context.Context,
	lang string,
) (*ReservationResponse, error) {
	res := &ReservationResponse{Rooms: []model.SpaceRoomCard{}}

	roomEntities, err := u.roomRepo.FindActiveRoomList(ctx, lang)
	if err != nil {
		return res, nil
	}

	res.Rooms = converter.SpaceRoomListToCardModel(
		roomEntities,
		lang,
		u.baseURL,
	)

	return res, nil
}
