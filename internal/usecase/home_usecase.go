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
	clientRepo  repository.ClientRepository
	newsRepo    repository.NewsArticleRepository

	baseURL string
}

type HomeResponse struct {
	Hero            []model.HeroSlide  `json:"hero"`
	Promo           []model.PromoSlide `json:"promo"`
	PartnerScroller []model.Partner    `json:"partner_scroller"`
	ClientScroller  []model.Client     `json:"client_scroller"`
	News            []model.NewsCard   `json:"news"`
}

func NewHomeUsecase(
	heroRepo repository.HeroRepository,
	promoRepo repository.PromoRepository,
	partnerRepo repository.PartnerRepository,
	clientRepo repository.ClientRepository,
	newsRepo repository.NewsArticleRepository,
	baseURL string,
) HomeUsecase {
	return &homeUsecaseImpl{
		heroRepo:    heroRepo,
		promoRepo:   promoRepo,
		partnerRepo: partnerRepo,
		clientRepo:  clientRepo,
		newsRepo:    newsRepo,
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

	clientEntities, err := u.clientRepo.FindActiveForScroller(ctx, lang)
	if err != nil {
		return nil, err
	}
	clientScroller := converter.ClientListToModel(
		clientEntities,
		lang,
		u.baseURL,
	)

	articles, err := u.newsRepo.FindLatest(ctx, lang, 5)
	if err != nil {
		return nil, err
	}

	news := make([]model.NewsCard, 0, len(articles))
	for _, a := range articles {
		card := converter.NewsArticleToNewsCard(a, lang, u.baseURL)
		if card != nil {
			news = append(news, *card)
		}
	}

	return &HomeResponse{
		Hero:            hero,
		Promo:           promo,
		PartnerScroller: partnerScroller,
		ClientScroller:  clientScroller,
		News:            news,
	}, nil
}
