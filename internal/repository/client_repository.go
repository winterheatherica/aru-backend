package repository

import (
	"context"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ClientRepository interface {
	FindActiveForScroller(ctx context.Context, lang string) ([]entity.Client, error)
	FindAll(ctx context.Context) ([]entity.Client, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error)
	Create(ctx context.Context, item *entity.Client) error
	Update(ctx context.Context, item *entity.Client) error
	UpsertTranslation(ctx context.Context, tr *entity.ClientTranslation) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type clientRepositoryImpl struct {
	db *gorm.DB
}

func NewClientRepository(db *gorm.DB) ClientRepository {
	return &clientRepositoryImpl{
		db: db,
	}
}

func (r *clientRepositoryImpl) FindActiveForScroller(
	ctx context.Context,
	lang string,
) ([]entity.Client, error) {

	var clients []entity.Client

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Where("is_active_client_scroller = ?", true).
		Order("order_index ASC").
		Find(&clients).Error

	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (r *clientRepositoryImpl) FindAll(ctx context.Context) ([]entity.Client, error) {
	var items []entity.Client
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Order("is_active_client_scroller DESC, order_index ASC, created_at DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *clientRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.Client, error) {
	var item entity.Client
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *clientRepositoryImpl) Create(ctx context.Context, item *entity.Client) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *clientRepositoryImpl) Update(ctx context.Context, item *entity.Client) error {
	return r.db.WithContext(ctx).
		Model(&entity.Client{}).
		Where("id = ?", item.ID).
		Updates(map[string]any{
			"image_path":                item.ImagePath,
			"order_index":               item.OrderIndex,
			"is_active_client_scroller": item.IsActiveClientScroller,
			"uploaded_by":               item.UploadedBy,
		}).Error
}

func (r *clientRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.ClientTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "client_id"}, {Name: "language"}},
			DoUpdates: clause.Assignments(map[string]any{
				"alt":         tr.Alt,
				"title":       tr.Title,
				"description": tr.Description,
			}),
		}).
		Create(tr).Error
}

func (r *clientRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.Client{}, "id = ?", id).Error
}
