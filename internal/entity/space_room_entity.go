package entity

import (
	"time"

	"github.com/google/uuid"
)

type SpaceRoom struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Code     string  `gorm:"column:code;not null" json:"code"`
	Capacity *int    `gorm:"column:capacity" json:"capacity,omitempty"`
	Floor    *int    `gorm:"column:floor" json:"floor,omitempty"`
	Address  *string `gorm:"column:address" json:"address,omitempty"`

	MainImageURL *string `gorm:"column:main_image_url" json:"main_image_url,omitempty"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []SpaceRoomTranslation `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"translations"`
	Images       []SpaceRoomImage       `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"images"`
	Bookings     []SpaceRoomBooking     `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"bookings"`
}

func (SpaceRoom) TableName() string {
	return "space_rooms"
}

type SpaceRoomTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	RoomID uuid.UUID `gorm:"column:room_id;type:uuid;not null" json:"room_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title       string   `gorm:"column:title;not null" json:"title"`
	Description *string  `gorm:"column:description" json:"description,omitempty"`
	Facilities  []string `gorm:"column:facilities;type:text[]" json:"facilities,omitempty"`

	MainImageAlt   *string `gorm:"column:main_image_alt" json:"main_image_alt,omitempty"`
	MainImageTitle *string `gorm:"column:main_image_title" json:"main_image_title,omitempty"`

	Slug string `gorm:"column:slug;not null" json:"slug"`

	MetaTitle       *string  `gorm:"column:meta_title" json:"meta_title,omitempty"`
	MetaDescription *string  `gorm:"column:meta_description" json:"meta_description,omitempty"`
	MetaKeywords    []string `gorm:"column:meta_keywords;type:text[]" json:"meta_keywords,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SpaceRoomTranslation) TableName() string {
	return "space_room_translations"
}

type SpaceRoomImage struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	RoomID uuid.UUID `gorm:"column:room_id;type:uuid;not null" json:"room_id"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	ImageURL string `gorm:"column:image_url;not null" json:"image_url"`

	IsActive   bool `gorm:"column:is_active;default:true" json:"is_active"`
	OrderIndex int  `gorm:"column:order_index;default:1" json:"order_index"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []SpaceRoomImageTranslation `gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (SpaceRoomImage) TableName() string {
	return "space_room_images"
}

type SpaceRoomImageTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ImageID uuid.UUID `gorm:"column:image_id;type:uuid;not null" json:"image_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Alt   *string `gorm:"column:alt" json:"alt,omitempty"`
	Title *string `gorm:"column:title" json:"title,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SpaceRoomImageTranslation) TableName() string {
	return "space_room_image_translations"
}

type SpaceRoomBooking struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	RoomID uuid.UUID `gorm:"column:room_id;type:uuid;not null" json:"room_id"`

	StartTime time.Time `gorm:"column:start_time;not null" json:"start_time"`
	EndTime   time.Time `gorm:"column:end_time;not null" json:"end_time"`

	Status string `gorm:"column:status;type:room_booking_status;default:PENDING" json:"status"`

	BookedBy *uuid.UUID `gorm:"column:booked_by" json:"booked_by"`
	AddedBy  *uuid.UUID `gorm:"column:added_by" json:"added_by"`

	Purpose *string `gorm:"column:purpose" json:"purpose,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SpaceRoomBooking) TableName() string {
	return "space_room_bookings"
}
