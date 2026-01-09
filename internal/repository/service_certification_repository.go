package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type ServiceCertificationRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServiceCertification, error)
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
