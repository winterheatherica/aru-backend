package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type SpaceRoomRepository interface {
	FindActiveRoomList(ctx context.Context, lang string) ([]entity.SpaceRoom, error)
	FindActiveByID(ctx context.Context, id string, lang string) (*entity.SpaceRoom, error)
	ResolveIDBySlug(ctx context.Context, slug string) (string, error)
	FindSlugByIDAndLang(ctx context.Context, id, lang string) (string, error)
}

type spaceRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewSpaceRoomRepository(db *gorm.DB) SpaceRoomRepository {
	return &spaceRoomRepositoryImpl{
		db: db,
	}
}

func (r *spaceRoomRepositoryImpl) FindActiveRoomList(
	ctx context.Context,
	lang string,
) ([]entity.SpaceRoom, error) {

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

func (r *spaceRoomRepositoryImpl) FindActiveByID(
	ctx context.Context,
	id string,
	lang string,
) (*entity.SpaceRoom, error) {

	var room entity.SpaceRoom

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Where("id = ?", id).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		Preload("Bookings").
		First(&room).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &room, nil
}

func (r *spaceRoomRepositoryImpl) ResolveIDBySlug(
	ctx context.Context,
	slug string,
) (string, error) {

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

func (r *spaceRoomRepositoryImpl) FindSlugByIDAndLang(
	ctx context.Context,
	id string,
	lang string,
) (string, error) {

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
