package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type HomeUsecase interface {
	GetHome(ctx context.Context, lang string) (*HomeResponse, error)
}

type homeUsecaseImpl struct {
	heroRepo    repository.HeroRepository
	promoRepo   repository.PromoRepository
	partnerRepo repository.PartnerRepository

	baseURL string
}

type HomeResponse struct {
	Hero            []model.HeroSlide  `json:"hero"`
	Promo           []model.PromoSlide `json:"promo"`
	PartnerScroller []model.Partner    `json:"partner_scroller"`
}

func NewHomeUsecase(
	heroRepo repository.HeroRepository,
	promoRepo repository.PromoRepository,
	partnerRepo repository.PartnerRepository,
	baseURL string,
) HomeUsecase {
	return &homeUsecaseImpl{
		heroRepo:    heroRepo,
		promoRepo:   promoRepo,
		partnerRepo: partnerRepo,
		baseURL:     baseURL,
	}
}

func (u *homeUsecaseImpl) GetHome(
	ctx context.Context,
	lang string,
) (*HomeResponse, error) {

	heroSlides, err := u.heroRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}
	hero := converter.HeroSlideListToModel(heroSlides, lang, u.baseURL)

	promoSlides, err := u.promoRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}
	promo := converter.PromoSlideListToModel(promoSlides, lang, u.baseURL)

	partnerEntities, err := u.partnerRepo.FindActiveForScroller(ctx, lang)
	if err != nil {
		return nil, err
	}
	partnerScroller := converter.PartnerListToModel(
		partnerEntities,
		lang,
		u.baseURL,
	)

	return &HomeResponse{
		Hero:            hero,
		Promo:           promo,
		PartnerScroller: partnerScroller,
	}, nil
}
