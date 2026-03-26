package repository

import (
	"context"
	"time"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CareerVacancyRepository interface {
	FindActiveOpen(ctx context.Context, lang string) ([]entity.CareerVacancy, error)
	FindAll(ctx context.Context) ([]entity.CareerVacancy, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CareerVacancy, error)
	Create(ctx context.Context, item *entity.CareerVacancy) error
	Update(ctx context.Context, item *entity.CareerVacancy) error
	UpsertTranslation(ctx context.Context, tr *entity.CareerVacancyTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *careerVacancyRepositoryImpl) FindAll(ctx context.Context) ([]entity.CareerVacancy, error) {
	var items []entity.CareerVacancy
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active DESC, opened_at DESC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *careerVacancyRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.CareerVacancy, error) {
	var item entity.CareerVacancy
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *careerVacancyRepositoryImpl) Create(ctx context.Context, item *entity.CareerVacancy) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *careerVacancyRepositoryImpl) Update(ctx context.Context, item *entity.CareerVacancy) error {
	return r.db.WithContext(ctx).
		Model(&entity.CareerVacancy{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"title":      item.Title,
			"slug":       item.Slug,
			"employment": item.Employment,
			"location":   item.Location,
			"opened_at":  item.OpenedAt,
			"closed_at":  item.ClosedAt,
			"is_active":  item.IsActive,
		}).Error
}

func (r *careerVacancyRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.CareerVacancyTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "vacancy_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"description": tr.Description,
			}),
		}).
		Create(tr).Error
}

func (r *careerVacancyRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.CareerVacancy{}, "id = ?", id).Error
}
