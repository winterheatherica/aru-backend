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
	partnerRepo repository.PartnerRepository
	awardRepo   repository.AwardRepository
	baseURL     string
}

type AboutResponse struct {
	Histories   []model.History `json:"histories"`
	PartnerGrid []model.Partner `json:"partner_grid"`
	Awards      []model.Award   `json:"awards"`
}

func NewAboutUsecase(
	historyRepo repository.HistoryRepository,
	partnerRepo repository.PartnerRepository,
	awardRepo repository.AwardRepository,
	baseURL string,
) AboutUsecase {
	return &aboutUsecaseImpl{
		historyRepo: historyRepo,
		partnerRepo: partnerRepo,
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

	partnerEntities, err := u.partnerRepo.FindActiveForGrid(ctx, lang)
	if err != nil {
		return nil, err
	}
	partnerGrid := converter.PartnerListToModel(
		partnerEntities,
		lang,
		u.baseURL,
	)

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
		Histories:   histories,
		PartnerGrid: partnerGrid,
		Awards:      awards,
	}, nil
}
