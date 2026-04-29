package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type InformationUsecase interface {
	GetInformation(ctx context.Context, lang string, year *int, page int) (*InformationResponse, error)
}

type informationUsecaseImpl struct {
	newsRepo repository.NewsArticleRepository
	baseURL  string
}

func NewInformationUsecase(
	newsRepo repository.NewsArticleRepository,
	baseURL string,
) InformationUsecase {
	return &informationUsecaseImpl{
		newsRepo: newsRepo,
		baseURL:  baseURL,
	}
}

type InformationResponse struct {
	NewsCards []model.NewsCard `json:"news_cards"`
	NewsYears []int            `json:"news_years"`
}

func (u *informationUsecaseImpl) GetInformation(
	ctx context.Context,
	lang string,
	year *int,
	page int,
) (*InformationResponse, error) {
	res := &InformationResponse{
		NewsCards: []model.NewsCard{},
		NewsYears: []int{},
	}

	_ = page

	if entities, err := u.newsRepo.FindActiveCardList(ctx, lang, year, 0, 0); err == nil {
		cards := make([]model.NewsCard, 0, len(entities))
		for _, e := range entities {
			card := converter.NewsArticleToNewsCard(e, lang, u.baseURL)
			if card != nil {
				cards = append(cards, *card)
			}
		}
		res.NewsCards = cards
	}

	if years, err := u.newsRepo.FindActiveYears(ctx); err == nil {
		res.NewsYears = years
	}

	return res, nil
}
