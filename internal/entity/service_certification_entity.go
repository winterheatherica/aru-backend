package entity

import (
	"time"

	"github.com/google/uuid"
)

type ServiceCertification struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Service string `gorm:"column:service;type:service;not null" json:"service"`

	OrderIndex int  `gorm:"column:order_index;default:1" json:"order_index"`
	IsActive   bool `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceCertificationTranslation `gorm:"foreignKey:CertificationID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (ServiceCertification) TableName() string {
	return "service_certifications"
}

type ServiceCertificationTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CertificationID uuid.UUID `gorm:"column:certification_id;type:uuid;not null" json:"certification_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title   *string `gorm:"column:title" json:"title,omitempty"`
	Alt     *string `gorm:"column:alt" json:"alt,omitempty"`
	Caption *string `gorm:"column:caption" json:"caption,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceCertificationTranslation) TableName() string {
	return "service_certification_translations"
}
