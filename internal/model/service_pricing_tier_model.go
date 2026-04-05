package model

import "github.com/google/uuid"

type ServicePricingTier struct {
	ID uuid.UUID `json:"id"`

	Service string `json:"service"`

	PriceMonthly float64 `json:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly"`

	Popular    bool `json:"popular"`
	OrderIndex int  `json:"order_index"`

	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Features    []string `json:"features"`

	IsActive bool `json:"is_active"`
}

type ServicePricingTierUpsertInput struct {
	Service      string   `json:"service"`
	Language     string   `json:"language"`
	PriceMonthly float64  `json:"price_monthly"`
	PriceYearly  float64  `json:"price_yearly"`
	Popular      bool     `json:"popular"`
	OrderIndex   int      `json:"order_index"`
	IsActive     bool     `json:"is_active"`
	Name         string   `json:"name"`
	Description  *string  `json:"description"`
	Features     []string `json:"features"`
}

type ServicePricingTierAdminTranslation struct {
	Language    string   `json:"language"`
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Features    []string `json:"features"`
}

type ServicePricingTierAdminItem struct {
	ID uuid.UUID `json:"id"`

	Service      string  `json:"service"`
	PriceMonthly float64 `json:"price_monthly"`
	PriceYearly  float64 `json:"price_yearly"`
	Popular      bool    `json:"popular"`
	OrderIndex   int     `json:"order_index"`
	IsActive     bool    `json:"is_active"`

	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Features    []string `json:"features"`

	Translations []ServicePricingTierAdminTranslation `json:"translations"`
}
