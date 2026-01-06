package entity

import (
	"time"

	"github.com/google/uuid"
)

type ServiceGallery struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Service   string `gorm:"column:service;type:service;not null" json:"service"`
	MediaType string `gorm:"column:media_type;type:gallery_media;not null" json:"media_type"`

	Src       string  `gorm:"column:src;not null" json:"src"`
	Thumbnail *string `gorm:"column:thumbnail" json:"thumbnail,omitempty"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceGalleryTranslation `gorm:"foreignKey:GalleryID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (ServiceGallery) TableName() string { return "service_galleries" }

type ServiceGalleryTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	GalleryID uuid.UUID `gorm:"column:gallery_id;type:uuid;not null" json:"gallery_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title   *string `gorm:"column:title" json:"title,omitempty"`
	Alt     *string `gorm:"column:alt" json:"alt,omitempty"`
	Caption *string `gorm:"column:caption" json:"caption,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceGalleryTranslation) TableName() string {
	return "service_gallery_translations"
}
