package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ServiceCertificationRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServiceCertification, error)
	FindByService(ctx context.Context, service string) ([]entity.ServiceCertification, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceCertification, error)
	Create(ctx context.Context, item *entity.ServiceCertification) error
	Update(ctx context.Context, item *entity.ServiceCertification) error
	UpsertTranslation(ctx context.Context, tr *entity.ServiceCertificationTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type serviceCertificationRepositoryImpl struct {
	db *gorm.DB
}

func NewServiceCertificationRepository(db *gorm.DB) ServiceCertificationRepository {
	return &serviceCertificationRepositoryImpl{db: db}
}

func (r *serviceCertificationRepositoryImpl) FindActiveByService(
	ctx context.Context,
	service string,
	lang string,
) ([]entity.ServiceCertification, error) {

	var certs []entity.ServiceCertification

	err := r.db.WithContext(ctx).
		Model(&entity.ServiceCertification{}).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("service = ?", service).
		Order("order_index ASC").
		Find(&certs).Error

	if err != nil {
		return nil, err
	}

	return certs, nil
}

func (r *serviceCertificationRepositoryImpl) FindByService(ctx context.Context, service string) ([]entity.ServiceCertification, error) {
	var certs []entity.ServiceCertification
	err := r.db.WithContext(ctx).
		Preload("Translations").
		Where("service = ?", service).
		Order("order_index ASC, created_at DESC").
		Find(&certs).Error
	if err != nil {
		return nil, err
	}
	return certs, nil
}

func (r *serviceCertificationRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceCertification, error) {
	var cert entity.ServiceCertification
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&cert, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *serviceCertificationRepositoryImpl) Create(ctx context.Context, item *entity.ServiceCertification) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *serviceCertificationRepositoryImpl) Update(ctx context.Context, item *entity.ServiceCertification) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *serviceCertificationRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.ServiceCertificationTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "certification_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "alt", "caption", "updated_at"}),
		}).
		Create(tr).Error
}

func (r *serviceCertificationRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ServiceCertification{}, "id = ?", id).Error
}
