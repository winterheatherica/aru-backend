package repository

import (
	"context"

	"aru-backend/internal/entity"

	"gorm.io/gorm"
)

type SpaceRoomRepository interface {
	FindActiveRooms(ctx context.Context, lang string) ([]entity.SpaceRoom, error)
	FindActiveBySlug(ctx context.Context, slug string, lang string) (*entity.SpaceRoom, error)
}

type spaceRoomRepositoryImpl struct {
	db *gorm.DB
}

func NewSpaceRoomRepository(db *gorm.DB) SpaceRoomRepository {
	return &spaceRoomRepositoryImpl{db: db}
}

func (r *spaceRoomRepositoryImpl) FindActiveRooms(
	ctx context.Context,
	lang string,
) ([]entity.SpaceRoom, error) {

	var rooms []entity.SpaceRoom

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		Where("is_active = ?", true).
		Order("created_at DESC").
		Find(&rooms).Error

	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *spaceRoomRepositoryImpl) FindActiveBySlug(
	ctx context.Context,
	slug string,
	lang string,
) (*entity.SpaceRoom, error) {

	var room entity.SpaceRoom

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		Preload("Bookings").
		Where("is_active = ?", true).
		Where("id IN (?)", r.subQueryRoomIDBySlug(slug, lang)).
		First(&room).Error

	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *spaceRoomRepositoryImpl) subQueryRoomIDBySlug(
	slug string,
	lang string,
) *gorm.DB {

	return r.db.
		Table("space_room_translations").
		Select("room_id").
		Where("slug = ?", slug).
		Where("language = ?", lang)
}

func (r *spaceRoomRepositoryImpl) FindActiveDetailBySlug(
	ctx context.Context,
	slug string,
	lang string,
) (*entity.SpaceRoom, error) {

	var room entity.SpaceRoom

	err := r.db.WithContext(ctx).
		Preload("Translations", "language = ?", lang).
		Preload("Images", "is_active = ?", true).
		Preload("Images.Translations", "language = ?", lang).
		Preload("Bookings").
		Where("is_active = ?", true).
		Where("id IN (?)", r.subQueryRoomIDBySlug(slug, lang)).
		First(&room).Error

	if err != nil {
		return nil, err
	}

	return &room, nil
}
