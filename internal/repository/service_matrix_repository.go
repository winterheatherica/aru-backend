package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceMatrixRepository interface {
	FindActiveByService(ctx context.Context, service string, compact *bool, lang string) (*entity.ServiceMatrix, error)
	FindByService(ctx context.Context, service string) ([]entity.ServiceMatrix, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceMatrix, error)
	FindByIDFull(ctx context.Context, id uuid.UUID) (*entity.ServiceMatrix, error)
	Create(ctx context.Context, item *entity.ServiceMatrix) error
	Update(ctx context.Context, item *entity.ServiceMatrix) error
	UpsertTranslation(ctx context.Context, tr *entity.ServiceMatrixTranslation) error
	ClearStructure(ctx context.Context, matrixID uuid.UUID) error
	CreateColumns(ctx context.Context, items []entity.ServiceMatrixColumn) error
	CreateColumnTranslations(ctx context.Context, items []entity.ServiceMatrixColumnTranslation) error
	CreateRows(ctx context.Context, items []entity.ServiceMatrixRow) error
	CreateRowTranslations(ctx context.Context, items []entity.ServiceMatrixRowTranslation) error
	CreateCells(ctx context.Context, items []entity.ServiceMatrixCell) error
	CreateCellTranslations(ctx context.Context, items []entity.ServiceMatrixCellTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *serviceMatrixRepositoryImpl) FindByService(ctx context.Context, service string) ([]entity.ServiceMatrix, error) {
	var items []entity.ServiceMatrix
	err := r.db.WithContext(ctx).
		Model(&entity.ServiceMatrix{}).
		Where("service = ?", service).
		Preload("Translations").
		Preload("Columns").
		Preload("Rows").
		Order("updated_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *serviceMatrixRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceMatrix, error) {
	var item entity.ServiceMatrix
	err := r.db.WithContext(ctx).
		Model(&entity.ServiceMatrix{}).
		Where("id = ?", id).
		Preload("Translations").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *serviceMatrixRepositoryImpl) FindByIDFull(ctx context.Context, id uuid.UUID) (*entity.ServiceMatrix, error) {
	var item entity.ServiceMatrix
	err := r.db.WithContext(ctx).
		Model(&entity.ServiceMatrix{}).
		Where("id = ?", id).
		Preload("Translations").
		Preload("Columns.Translations").
		Preload("Rows.Translations").
		Preload("Rows.Cells.Translations").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *serviceMatrixRepositoryImpl) Create(ctx context.Context, item *entity.ServiceMatrix) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *serviceMatrixRepositoryImpl) Update(ctx context.Context, item *entity.ServiceMatrix) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *serviceMatrixRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.ServiceMatrixTranslation) error {
	return r.db.WithContext(ctx).
		Where("matrix_id = ? AND language = ?", tr.MatrixID, tr.Language).
		Assign(map[string]any{
			"title":       tr.Title,
			"description": tr.Description,
			"footnote":    tr.Footnote,
		}).
		FirstOrCreate(tr).Error
}

func (r *serviceMatrixRepositoryImpl) ClearStructure(ctx context.Context, matrixID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("matrix_id = ?", matrixID).Delete(&entity.ServiceMatrixRow{}).Error; err != nil {
			return err
		}
		if err := tx.Where("matrix_id = ?", matrixID).Delete(&entity.ServiceMatrixColumn{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *serviceMatrixRepositoryImpl) CreateColumns(ctx context.Context, items []entity.ServiceMatrixColumn) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) CreateColumnTranslations(ctx context.Context, items []entity.ServiceMatrixColumnTranslation) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) CreateRows(ctx context.Context, items []entity.ServiceMatrixRow) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) CreateRowTranslations(ctx context.Context, items []entity.ServiceMatrixRowTranslation) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) CreateCells(ctx context.Context, items []entity.ServiceMatrixCell) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) CreateCellTranslations(ctx context.Context, items []entity.ServiceMatrixCellTranslation) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *serviceMatrixRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ServiceMatrix{}, "id = ?", id).Error
}
