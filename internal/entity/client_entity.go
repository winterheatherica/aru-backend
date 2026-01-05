package entity

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ImagePath string `gorm:"column:image_path;not null" json:"image_path"`

	OrderIndex int `gorm:"column:order_index;default:0" json:"order_index"`

	IsActiveClientScroller bool `gorm:"column:is_active_client_scroller;default:true" json:"is_active_client_scroller"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ClientTranslation `gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (Client) TableName() string {
	return "clients"
}

type ClientTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ClientID uuid.UUID `gorm:"column:client_id;type:uuid;not null" json:"client_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title       *string `gorm:"column:title" json:"title,omitempty"`
	Alt         *string `gorm:"column:alt" json:"alt,omitempty"`
	Description *string `gorm:"column:description" json:"description,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ClientTranslation) TableName() string {
	return "client_translations"
}
