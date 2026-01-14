package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type ArticleUsecase interface {
	GetArticleByID(ctx context.Context, id string, lang string) (*model.NewsArticle, error)
}

type articleUsecaseImpl struct {
	articleRepo repository.NewsArticleRepository
	baseURL     string
}

func NewArticleUsecase(
	articleRepo repository.NewsArticleRepository,
	baseURL string,
) ArticleUsecase {
	return &articleUsecaseImpl{
		articleRepo: articleRepo,
		baseURL:     baseURL,
	}
}

func (u *articleUsecaseImpl) GetArticleByID(
	ctx context.Context,
	id string,
	lang string,
) (*model.NewsArticle, error) {

	entity, err := u.articleRepo.FindActiveByID(ctx, id, lang)
	if err != nil {
		return nil, err
	}

	return converter.NewsArticleToModel(*entity, lang, u.baseURL), nil
}
