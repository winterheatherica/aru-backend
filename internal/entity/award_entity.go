package entity

import (
	"time"

	"github.com/google/uuid"
)

type Award struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ImagePath string `gorm:"column:image_path;not null" json:"image_path"`

	Year       *int `gorm:"column:year" json:"year,omitempty"`
	OrderIndex int  `gorm:"column:order_index;default:0" json:"order_index"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []AwardTranslation `gorm:"foreignKey:AwardID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (Award) TableName() string {
	return "awards"
}

type AwardTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	AwardID uuid.UUID `gorm:"column:award_id;type:uuid;not null" json:"award_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Alt         *string `gorm:"column:alt" json:"alt,omitempty"`
	Title       *string `gorm:"column:title" json:"title,omitempty"`
	Label       *string `gorm:"column:label" json:"label,omitempty"`
	Description *string `gorm:"column:description" json:"description,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (AwardTranslation) TableName() string {
	return "award_translations"
}
