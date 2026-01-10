package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ServicePricingTier struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Service string `gorm:"column:service;type:service;not null" json:"service"`

	PriceMonthly float64 `gorm:"column:price_monthly;type:numeric(12,2);not null" json:"price_monthly"`
	PriceYearly  float64 `gorm:"column:price_yearly;type:numeric(12,2);not null" json:"price_yearly"`

	Popular    bool `gorm:"column:popular;default:false" json:"popular"`
	OrderIndex int  `gorm:"column:order_index;default:1" json:"order_index"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServicePricingTierTranslation `gorm:"foreignKey:TierID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (ServicePricingTier) TableName() string { return "service_pricing_tiers" }

type ServicePricingTierTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	TierID uuid.UUID `gorm:"column:tier_id;type:uuid;not null" json:"tier_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Name        string         `gorm:"column:name;not null" json:"name"`
	Description *string        `gorm:"column:description" json:"description,omitempty"`
	Features    pq.StringArray `gorm:"column:features;type:text[]" json:"features"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServicePricingTierTranslation) TableName() string {
	return "service_pricing_tier_translations"
}
