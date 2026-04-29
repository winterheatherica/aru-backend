package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type CareerUsecase interface {
	GetCareers(ctx context.Context, lang string) ([]model.CareerVacancy, error)
}

type careerUsecaseImpl struct {
	vacancyRepo repository.CareerVacancyRepository
}

func NewCareerUsecase(
	vacancyRepo repository.CareerVacancyRepository,
) CareerUsecase {
	return &careerUsecaseImpl{
		vacancyRepo: vacancyRepo,
	}
}

func (u *careerUsecaseImpl) GetCareers(
	ctx context.Context,
	lang string,
) ([]model.CareerVacancy, error) {
	entities, err := u.vacancyRepo.FindActiveOpen(ctx, lang)
	if err != nil {
		return []model.CareerVacancy{}, nil
	}

	return converter.CareerVacancyListToModel(entities, lang), nil
}
