package entity

import (
	"time"

	"github.com/google/uuid"
)

type PromoSlide struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ImagePath  string    `gorm:"column:image_path;not null" json:"image_path"`
	OrderIndex int       `gorm:"column:order_index;default:0" json:"order_index"`
	IsActive   bool      `gorm:"column:is_active;default:true" json:"is_active"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []PromoSlideTranslation `gorm:"foreignKey:PromoSlideID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (PromoSlide) TableName() string {
	return "promo_slides"
}

type PromoSlideTranslation struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PromoSlideID uuid.UUID `gorm:"column:promo_slide_id;type:uuid;not null" json:"promo_slide_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title *string `gorm:"column:title" json:"title,omitempty"`
	Alt   *string `gorm:"column:alt" json:"alt,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PromoSlideTranslation) TableName() string {
	return "promo_slide_translations"
}
