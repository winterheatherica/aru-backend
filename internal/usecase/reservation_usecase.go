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

	roomEntities, err := u.roomRepo.FindActiveRoomList(ctx, lang)
	if err != nil {
		return nil, err
	}

	rooms := converter.SpaceRoomListToCardModel(
		roomEntities,
		lang,
		u.baseURL,
	)

	return &ReservationResponse{
		Rooms: rooms,
	}, nil
}
