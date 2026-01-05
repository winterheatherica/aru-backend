package repository

import (
	"context"
	"time"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type CareerVacancyRepository interface {
	FindActiveOpen(ctx context.Context, lang string) ([]entity.CareerVacancy, error)
}

type careerVacancyRepositoryImpl struct {
	db *gorm.DB
}

func NewCareerVacancyRepository(db *gorm.DB) CareerVacancyRepository {
	return &careerVacancyRepositoryImpl{
		db: db,
	}
}

func (r *careerVacancyRepositoryImpl) FindActiveOpen(
	ctx context.Context,
	lang string,
) ([]entity.CareerVacancy, error) {

	var vacancies []entity.CareerVacancy
	now := time.Now()

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("opened_at <= ?", now).
		Where("(closed_at IS NULL OR closed_at > ?)", now).
		Order("opened_at DESC").
		Find(&vacancies).Error

	if err != nil {
		return nil, err
	}

	return vacancies, nil
}
