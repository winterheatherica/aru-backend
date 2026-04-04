package repository

import (
	"context"
	"strings"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpaceRoomRepository interface {
	FindActiveRoomList(ctx context.Context, lang string) ([]entity.SpaceRoom, error)
	FindActiveByID(ctx context.Context, id string, lang string) (*entity.SpaceRoom, error)
	ResolveIDBySlug(ctx context.Context, slug string) (string, error)
	FindSlugByIDAndLang(ctx context.Context, id, lang string) (string, error)

	FindAll(ctx context.Context) ([]entity.SpaceRoom, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.SpaceRoom, error)
	Create(ctx context.Context, room *entity.SpaceRoom) error
	Update(ctx context.Context, room *entity.SpaceRoom) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
	UpsertTranslation(ctx context.Context, tr *entity.SpaceRoomTranslation) error
	CreateImage(ctx context.Context, image *entity.SpaceRoomImage) error
	UpsertImageTranslation(ctx context.Context, tr *entity.SpaceRoomImageTranslation) error
	DeleteImagesByRoomID(ctx context.Context, roomID uuid.UUID) error
	IsSlugExists(ctx context.Context, slug, lang string, excludeRoomID *uuid.UUID) (bool, error)
	FindUploaderByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.User, error)
}

type spaceRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewSpaceRoomRepository(db *gorm.DB) SpaceRoomRepository {
	return &spaceRoomRepositoryImpl{db: db}
}

func (r *spaceRoomRepositoryImpl) FindActiveRoomList(ctx context.Context, lang string) ([]entity.SpaceRoom, error) {
	var rooms []entity.SpaceRoom
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		Order("created_at DESC").
		Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *spaceRoomRepositoryImpl) FindActiveByID(ctx context.Context, id string, lang string) (*entity.SpaceRoom, error) {
	var room entity.SpaceRoom
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("id = ?", id).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *spaceRoomRepositoryImpl) ResolveIDBySlug(ctx context.Context, slug string) (string, error) {
	var id string
	err := r.db.WithContext(ctx).
		Table("space_room_translations").
		Select("room_id").
		Where("slug = ?", slug).
		Limit(1).
		Scan(&id).Error
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *spaceRoomRepositoryImpl) FindSlugByIDAndLang(ctx context.Context, id, lang string) (string, error) {
	var slug string
	err := r.db.WithContext(ctx).
		Table("space_room_translations").
		Select("slug").
		Where("room_id = ?", id).
		Where("language = ?", lang).
		Limit(1).
		Scan(&slug).Error
	if err != nil {
		return "", err
	}
	if slug == "" {
		return "", nil
	}
	return slug, nil
}

func (r *spaceRoomRepositoryImpl) FindAll(ctx context.Context) ([]entity.SpaceRoom, error) {
	var rooms []entity.SpaceRoom
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Images").
		Preload("Images.Translations").
		Order("created_at DESC").
		Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *spaceRoomRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.SpaceRoom, error) {
	var room entity.SpaceRoom
	if err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Images").
		Preload("Images.Translations").
		First(&room, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *spaceRoomRepositoryImpl) Create(ctx context.Context, room *entity.SpaceRoom) error {
	return r.db.WithContext(ctx).Create(room).Error
}

func (r *spaceRoomRepositoryImpl) Update(ctx context.Context, room *entity.SpaceRoom) error {
	return r.db.WithContext(ctx).Model(&entity.SpaceRoom{}).
		Where("id = ?", room.ID).
		Updates(map[string]any{
			"code":           room.Code,
			"capacity":       room.Capacity,
			"floor":          room.Floor,
			"address":        room.Address,
			"main_image_url": room.MainImageURL,
			"uploaded_by":    room.UploadedBy,
			"is_active":      room.IsActive,
		}).Error
}

func (r *spaceRoomRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SpaceRoom{}, "id = ?", id).Error
}

func (r *spaceRoomRepositoryImpl) UpsertTranslation(ctx context.Context, tr *entity.SpaceRoomTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "room_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "description", "facilities", "main_image_alt", "main_image_title", "slug", "meta_title", "meta_description", "meta_keywords", "updated_at"}),
		}).
		Create(tr).Error
}

func (r *spaceRoomRepositoryImpl) CreateImage(ctx context.Context, image *entity.SpaceRoomImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *spaceRoomRepositoryImpl) UpsertImageTranslation(ctx context.Context, tr *entity.SpaceRoomImageTranslation) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "image_id"}, {Name: "language"}},
			DoUpdates: clause.AssignmentColumns([]string{"alt", "title", "updated_at"}),
		}).Create(tr).Error
}

func (r *spaceRoomRepositoryImpl) DeleteImagesByRoomID(ctx context.Context, roomID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.SpaceRoomImage{}, "room_id = ?", roomID).Error
}

func (r *spaceRoomRepositoryImpl) IsSlugExists(ctx context.Context, slug, lang string, excludeRoomID *uuid.UUID) (bool, error) {
	q := r.db.WithContext(ctx).Table("space_room_translations").Where("slug = ?", slug).Where("language = ?", strings.ToUpper(lang))
	if excludeRoomID != nil {
		q = q.Where("room_id <> ?", *excludeRoomID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *spaceRoomRepositoryImpl) FindUploaderByRoomID(ctx context.Context, roomID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Joins("JOIN space_rooms sr ON sr.uploaded_by = users.id").
		Where("sr.id = ?", roomID).
		Limit(1).
		First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
