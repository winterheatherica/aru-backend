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
	res := &HomeResponse{
		Hero:            []model.HeroSlide{},
		Promo:           []model.PromoSlide{},
		PartnerScroller: []model.Partner{},
		ClientScroller:  []model.Client{},
		News:            []model.NewsCard{},
	}

	if heroSlides, err := u.heroRepo.FindActiveByLanguage(ctx, lang); err == nil {
		res.Hero = converter.HeroSlideListToModel(heroSlides, lang, u.baseURL)
	}

	if promoSlides, err := u.promoRepo.FindActiveByLanguage(ctx, lang); err == nil {
		res.Promo = converter.PromoSlideListToModel(promoSlides, lang, u.baseURL)
	}

	if partnerEntities, err := u.partnerRepo.FindActiveForScroller(ctx, lang); err == nil {
		res.PartnerScroller = converter.PartnerListToModel(
			partnerEntities,
			lang,
			u.baseURL,
		)
	}

	if clientEntities, err := u.clientRepo.FindActiveForScroller(ctx, lang); err == nil {
		res.ClientScroller = converter.ClientListToModel(
			clientEntities,
			lang,
			u.baseURL,
		)
	}

	if articles, err := u.newsRepo.FindLatest(ctx, lang, 5); err == nil {
		news := make([]model.NewsCard, 0, len(articles))
		for _, a := range articles {
			card := converter.NewsArticleToNewsCard(a, lang, u.baseURL)
			if card != nil {
				news = append(news, *card)
			}
		}
		res.News = news
	}

	return res, nil
}
