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
	heroRepo repository.HeroRepository
	baseURL  string
}

type HomeResponse struct {
	Hero []model.HeroSlide `json:"hero"`
}

func NewHomeUsecase(
	heroRepo repository.HeroRepository,
	baseURL string,
) HomeUsecase {
	return &homeUsecaseImpl{
		heroRepo: heroRepo,
		baseURL:  baseURL,
	}
}

func (u *homeUsecaseImpl) GetHome(ctx context.Context, lang string) (*HomeResponse, error) {
	slides, err := u.heroRepo.FindActiveByLanguage(ctx, lang)
	if err != nil {
		return nil, err
	}

	hero := converter.HeroSlideListToModel(slides, lang, u.baseURL)

	return &HomeResponse{
		Hero: hero,
	}, nil
}
