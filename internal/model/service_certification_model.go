package model

import "github.com/google/uuid"

type ServiceCertification struct {
	ID uuid.UUID `json:"id"`

	Service string `json:"service"`

	OrderIndex int  `json:"order_index"`
	IsActive   bool `json:"is_active"`

	Title   *string `json:"title,omitempty"`
	Alt     *string `json:"alt,omitempty"`
	Caption *string `json:"caption,omitempty"`
}

type ServiceCertificationAdminTranslation struct {
	Language string  `json:"language"`
	Title    *string `json:"title,omitempty"`
	Alt      *string `json:"alt,omitempty"`
	Caption  *string `json:"caption,omitempty"`
}

type ServiceCertificationAdminItem struct {
	ID          uuid.UUID                              `json:"id"`
	Service     string                                 `json:"service"`
	OrderIndex  int                                    `json:"order_index"`
	IsActive    bool                                   `json:"is_active"`
	Title       *string                                `json:"title,omitempty"`
	Alt         *string                                `json:"alt,omitempty"`
	Caption     *string                                `json:"caption,omitempty"`
	Translations []ServiceCertificationAdminTranslation `json:"translations"`
}

type ServiceCertificationUpsertInput struct {
	Service    string  `json:"service"`
	Language   string  `json:"language"`
	Title      *string `json:"title"`
	Alt        *string `json:"alt"`
	Caption    *string `json:"caption"`
	OrderIndex int     `json:"order_index"`
	IsActive   bool    `json:"is_active"`
}
