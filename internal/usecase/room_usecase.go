package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type RoomUsecase interface {
	GetRoomDetail(ctx context.Context, slug string, lang string) (*model.SpaceRoomDetail, error)
}

type roomUsecaseImpl struct {
	roomRepo repository.SpaceRoomRepository
	baseURL  string
}

func NewRoomUsecase(
	roomRepo repository.SpaceRoomRepository,
	baseURL string,
) RoomUsecase {
	return &roomUsecaseImpl{
		roomRepo: roomRepo,
		baseURL:  baseURL,
	}
}

func (u *roomUsecaseImpl) GetRoomDetail(
	ctx context.Context,
	slug string,
	lang string,
) (*model.SpaceRoomDetail, error) {

	roomEntity, err := u.roomRepo.FindActiveBySlug(ctx, slug, lang)
	if err != nil {
		return nil, err
	}

	detail := converter.SpaceRoomToDetailModel(
		*roomEntity,
		lang,
		u.baseURL,
	)

	return detail, nil
}
