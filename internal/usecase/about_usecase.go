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
	res := &AboutResponse{
		Histories:   []model.History{},
		PartnerGrid: []model.Partner{},
		Awards:      []model.Award{},
	}

	if historyEntities, err := u.historyRepo.FindActiveByLanguage(ctx, lang); err == nil {
		res.Histories = converter.HistoryListToModel(historyEntities)
	}

	if partnerEntities, err := u.partnerRepo.FindActiveForGrid(ctx, lang); err == nil {
		res.PartnerGrid = converter.PartnerListToModel(
			partnerEntities,
			lang,
			u.baseURL,
		)
	}

	if awardEntities, err := u.awardRepo.FindActiveByLanguage(ctx, lang); err == nil {
		res.Awards = converter.AwardListToModel(
			awardEntities,
			lang,
			u.baseURL,
		)
	}

	return res, nil
}
