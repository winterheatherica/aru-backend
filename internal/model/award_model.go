package model

import "github.com/google/uuid"

type Award struct {
	ID          uuid.UUID `json:"id"`
	Src         string    `json:"src"`
	Alt         string    `json:"alt"`
	Title       *string   `json:"title,omitempty"`
	Label       *string   `json:"label,omitempty"`
	Description *string   `json:"description,omitempty"`
	Year        *int      `json:"year,omitempty"`
	Order       int       `json:"order"`
}
