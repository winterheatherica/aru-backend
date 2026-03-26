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

type AwardUpsertInput struct {
	Language    string  `json:"language"`
	Alt         *string `json:"alt"`
	Title       *string `json:"title"`
	Label       *string `json:"label"`
	Description *string `json:"description"`
	Year        *int    `json:"year"`
	OrderIndex  int     `json:"order_index"`
	IsActive    bool    `json:"is_active"`
}

type AwardAdminTranslation struct {
	Language    string  `json:"language"`
	Alt         *string `json:"alt,omitempty"`
	Title       *string `json:"title,omitempty"`
	Label       *string `json:"label,omitempty"`
	Description *string `json:"description,omitempty"`
}

type AwardAdminItem struct {
	ID           uuid.UUID               `json:"id"`
	ImagePath    string                  `json:"image_path"`
	ImageURL     string                  `json:"image_url"`
	Year         *int                    `json:"year,omitempty"`
	OrderIndex   int                     `json:"order_index"`
	IsActive     bool                    `json:"is_active"`
	Translations []AwardAdminTranslation `json:"translations"`
}
