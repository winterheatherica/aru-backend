package model

import "github.com/google/uuid"

type HeroSlide struct {
	ID       uuid.UUID `json:"id"`
	Src      string    `json:"src"`
	Alt      string    `json:"alt"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Banner      string    `json:"banner"`
	Order    int       `json:"order"`
}

type HeroUpsertInput struct {
	Language    string  `json:"language"`
	Alt         *string `json:"alt"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Banner      string  `json:"banner"`
	OrderIndex int     `json:"order_index"`
	IsActive   bool    `json:"is_active"`
}

type HeroAdminTranslation struct {
	Language    string  `json:"language"`
	Alt         *string `json:"alt,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type HeroAdminItem struct {
	ID           uuid.UUID              `json:"id"`
	ImagePath    string                 `json:"image_path"`
	MainImageURL string                 `json:"main_image_url"`
	OrderIndex   int                    `json:"order_index"`
	IsActive     bool                   `json:"is_active"`
	Banner       string                 `json:"banner"`
	Translations []HeroAdminTranslation `json:"translations"`
}
