package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceGalleryRepository interface {
	FindActiveByService(ctx context.Context, service string, lang string) ([]entity.ServiceGallery, error)
	FindByService(ctx context.Context, service string) ([]entity.ServiceGallery, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceGallery, error)
	Create(ctx context.Context, item *entity.ServiceGallery) error
	Update(ctx context.Context, item *entity.ServiceGallery) error
	UpsertTranslation(ctx context.Context, tr *entity.ServiceGalleryTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
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

func (r *serviceGalleryRepositoryImpl) FindByService(ctx context.Context, service string) ([]entity.ServiceGallery, error) {
	var items []entity.ServiceGallery
	err := r.db.WithContext(ctx).
		Model(&entity.ServiceGallery{}).
		Where("service = ?", service).
		Preload("Translations").
		Order("updated_at DESC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *serviceGalleryRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.ServiceGallery, error) {
	var item entity.ServiceGallery
	err := r.db.WithContext(ctx).
		Model(&entity.ServiceGallery{}).
		Where("id = ?", id).
		Preload("Translations").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *serviceGalleryRepositoryImpl) Create(ctx context.Context, item *entity.ServiceGallery) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *serviceGalleryRepositoryImpl) Update(ctx context.Context, item *entity.ServiceGallery) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *serviceGalleryRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.ServiceGalleryTranslation) error {
	return r.db.WithContext(ctx).
		Where("gallery_id = ? AND language = ?", tr.GalleryID, tr.Language).
		Assign(map[string]any{
			"title":   tr.Title,
			"alt":     tr.Alt,
			"caption": tr.Caption,
		}).
		FirstOrCreate(tr).Error
}

func (r *serviceGalleryRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.ServiceGallery{}, "id = ?", id).Error
}
