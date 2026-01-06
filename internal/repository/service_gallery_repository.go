package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type ServiceGalleryRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServiceGallery, error)
}

type serviceGalleryRepositoryImpl struct {
	db *gorm.DB
}

func NewServiceGalleryRepository(db *gorm.DB) ServiceGalleryRepository {
	return &serviceGalleryRepositoryImpl{db: db}
}

func (r *serviceGalleryRepositoryImpl) FindActiveByService(
	ctx context.Context,
	service string,
	lang string,
) ([]entity.ServiceGallery, error) {

	var galleries []entity.ServiceGallery

	err := r.db.WithContext(ctx).
		Model(&entity.ServiceGallery{}).
		Preload("Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Where("service = ?", service).
		Order("created_at DESC").
		Limit(5).
		Find(&galleries).Error

	if err != nil {
		return nil, err
	}

	return galleries, nil
}
