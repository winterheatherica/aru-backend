package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type SiteContentUsecase interface {
	GetSiteContent(ctx context.Context, lang string) (*SiteContentResponse, error)
}

type siteContentUsecaseImpl struct {
	heroRepo repository.HeroRepository
	baseURL  string
}

type SiteContentResponse struct {
	Home HomeSection `json:"home"`
}

type HomeSection struct {
	Hero []model.HeroSlide `json:"hero"`
}

func NewSiteContentUsecase(heroRepo repository.HeroRepository, baseURL string) SiteContentUsecase {
	return &siteContentUsecaseImpl{
		heroRepo: heroRepo,
		baseURL:  baseURL,
	}
}

func (u *siteContentUsecaseImpl) GetSiteContent(ctx context.Context, lang string) (*SiteContentResponse, error) {
	slides, err := u.heroRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}

	hero := converter.HeroSlideListToModel(slides, lang, u.baseURL)

	home := HomeSection{
		Hero: hero,
	}

	return &SiteContentResponse{
		Home: home,
	}, nil
}
