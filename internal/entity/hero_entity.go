package entity

import (
	"time"

	"github.com/google/uuid"
)

type HeroSlide struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ImagePath  string    `gorm:"column:image_path;not null" json:"image_path"`
	OrderIndex int       `gorm:"column:order_index;default:0" json:"order_index"`
	IsActive   bool      `gorm:"column:is_active;default:true" json:"is_active"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []HeroSlideTranslation `gorm:"foreignKey:HeroSlideID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (HeroSlide) TableName() string {
	return "hero_slides"
}

type HeroSlideTranslation struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	HeroSlideID uuid.UUID `gorm:"column:hero_slide_id;type:uuid;not null" json:"hero_slide_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`
	Banner   string `gorm:"column:banner;type:hero;not null;default:POLISH" json:"banner"`

	Alt      *string `gorm:"column:alt" json:"alt,omitempty"`
	Title    *string `gorm:"column:title" json:"title,omitempty"`
	CtaLabel *string `gorm:"column:cta_label" json:"cta_label,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (HeroSlideTranslation) TableName() string {
	return "hero_slide_translations"
}
