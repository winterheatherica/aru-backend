package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type CategoryUsecase interface {
	GetCategoryBySlug(ctx context.Context, slug string, lang string) (*model.NewsCategory, error)
}

type categoryUsecaseImpl struct {
	categoryRepo repository.NewsCategoryRepository
}

func NewCategoryUsecase(
	categoryRepo repository.NewsCategoryRepository,
) CategoryUsecase {
	return &categoryUsecaseImpl{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecaseImpl) GetCategoryBySlug(
	ctx context.Context,
	slug string,
	lang string,
) (*model.NewsCategory, error) {

	entity, err := u.categoryRepo.FindActiveBySlug(ctx, slug, lang)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return nil, nil
	}

	category := converter.NewsCategoryToModel(*entity, lang)
	if category == nil {
		return nil, nil
	}

	return category, nil
}
