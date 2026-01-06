package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type ServiceMatrixRepository interface {
	FindActiveByService(ctx context.Context, service string, compact *bool, lang string) (*entity.ServiceMatrix, error)
}

type serviceMatrixRepositoryImpl struct {
	db *gorm.DB
}

func NewServiceMatrixRepository(db *gorm.DB) ServiceMatrixRepository {
	return &serviceMatrixRepositoryImpl{db: db}
}

func (r *serviceMatrixRepositoryImpl) FindActiveByService(
	ctx context.Context,
	service string,
	compact *bool,
	lang string,
) (*entity.ServiceMatrix, error) {

	var matrix entity.ServiceMatrix

	q := r.db.WithContext(ctx).
		Model(&entity.ServiceMatrix{}).
		Preload("Translations", "language = ?", lang).
		Preload("Columns.Translations", "language = ?", lang).
		Preload("Rows.Translations", "language = ?", lang).
		Preload("Rows.Cells.Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("service = ?", service)

	if compact != nil {
		q = q.Where("compact = ?", *compact)
	}

	err := q.First(&matrix).Error
	if err != nil {
		return nil, err
	}

	return &matrix, nil
}
