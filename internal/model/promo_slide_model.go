package model

import "github.com/google/uuid"

type PromoSlide struct {
	ID    uuid.UUID `json:"id"`
	Src   string    `json:"src"`
	Alt   string    `json:"alt"`
	Title *string   `json:"title,omitempty"`
	Order int       `json:"order"`
}
