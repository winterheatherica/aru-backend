package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type AboutUsecase interface {
	GetAbout(ctx context.Context, lang string) (*AboutResponse, error)
}

type aboutUsecaseImpl struct {
	historyRepo repository.HistoryRepository
}

type AboutResponse struct {
	Histories []model.History `json:"histories"`
}

func NewAboutUsecase(
	historyRepo repository.HistoryRepository,
) AboutUsecase {
	return &aboutUsecaseImpl{
		historyRepo: historyRepo,
	}
}

func (u *aboutUsecaseImpl) GetAbout(
	ctx context.Context,
	lang string,
) (*AboutResponse, error) {

	historyEntities, err := u.historyRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}

	histories := converter.HistoryListToModel(historyEntities)

	return &AboutResponse{
		Histories: histories,
	}, nil
}
