package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type RoomUsecase interface {
	GetRoomByID(ctx context.Context, id string, lang string) (*model.SpaceRoomDetail, error)
	ResolveRoomID(ctx context.Context, slug string) (string, error)
	ResolveRoomSlug(ctx context.Context, id string, lang string) (string, error)
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

func (u *roomUsecaseImpl) GetRoomByID(
	ctx context.Context,
	id string,
	lang string,
) (*model.SpaceRoomDetail, error) {

	entity, err := u.roomRepo.FindActiveByID(ctx, id, lang)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	return converter.SpaceRoomToDetailModel(*entity, lang, u.baseURL), nil
}

func (u *roomUsecaseImpl) ResolveRoomID(
	ctx context.Context,
	slug string,
) (string, error) {
	return u.roomRepo.ResolveIDBySlug(ctx, slug)
}

func (u *roomUsecaseImpl) ResolveRoomSlug(
	ctx context.Context,
	id string,
	lang string,
) (string, error) {
	return u.roomRepo.FindSlugByIDAndLang(ctx, id, lang)
}
