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
