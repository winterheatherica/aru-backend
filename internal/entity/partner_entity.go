package entity

import (
	"time"

	"github.com/google/uuid"
)

type Partner struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ImagePath string `gorm:"column:image_path;not null" json:"image_path"`

	OrderIndex int `gorm:"column:order_index;default:0" json:"order_index"`

	IsActivePartnerGrid     bool `gorm:"column:is_active_partner_grid;default:true" json:"is_active_partner_grid"`
	IsActivePartnerScroller bool `gorm:"column:is_active_partner_scroller;default:false" json:"is_active_partner_scroller"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []PartnerTranslation `gorm:"foreignKey:PartnerID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (Partner) TableName() string {
	return "partners"
}

type PartnerTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	PartnerID uuid.UUID `gorm:"column:partner_id;type:uuid;not null" json:"partner_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title       *string `gorm:"column:title" json:"title,omitempty"`
	Alt         *string `gorm:"column:alt" json:"alt,omitempty"`
	Description *string `gorm:"column:description" json:"description,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PartnerTranslation) TableName() string {
	return "partner_translations"
}
