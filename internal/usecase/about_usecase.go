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
	awardRepo   repository.AwardRepository
	baseURL     string
}

type AboutResponse struct {
	Histories []model.History `json:"histories"`
	Awards    []model.Award   `json:"awards"`
}

func NewAboutUsecase(
	historyRepo repository.HistoryRepository,
	awardRepo repository.AwardRepository,
	baseURL string,
) AboutUsecase {
	return &aboutUsecaseImpl{
		historyRepo: historyRepo,
		awardRepo:   awardRepo,
		baseURL:     baseURL,
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

	awardEntities, err := u.awardRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}
	awards := converter.AwardListToModel(
		awardEntities,
		lang,
		u.baseURL,
	)

	return &AboutResponse{
		Histories: histories,
		Awards:    awards,
	}, nil
}
