package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type InformationUsecase interface {
	GetInformationCards(ctx context.Context, lang string, year *int, page int, limit int) ([]model.NewsCard, error)
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

func (u *informationUsecaseImpl) GetInformationCards(
	ctx context.Context,
	lang string,
	year *int,
	page int,
	limit int,
) ([]model.NewsCard, error) {

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 18
	}

	offset := (page - 1) * limit

	entities, err := u.newsRepo.FindActiveCardList(ctx, lang, year, limit, offset)
	if err != nil {
		return nil, err
	}

	cards := make([]model.NewsCard, 0, len(entities))

	for _, e := range entities {
		card := converter.NewsArticleToNewsCard(e, lang, u.baseURL)
		if card != nil {
			cards = append(cards, *card)
		}
	}

	return cards, nil
}
