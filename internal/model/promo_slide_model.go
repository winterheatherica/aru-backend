package model

import "github.com/google/uuid"

type PromoSlide struct {
	ID    uuid.UUID `json:"id"`
	Src   string    `json:"src"`
	Alt   string    `json:"alt"`
	Title *string   `json:"title,omitempty"`
	Order int       `json:"order"`
}

type PromoSlideUpsertInput struct {
	Language   string  `json:"language"`
	Alt        *string `json:"alt"`
	Title      *string `json:"title"`
	OrderIndex int     `json:"order_index"`
	IsActive   bool    `json:"is_active"`
}

type PromoSlideAdminTranslation struct {
	Language string  `json:"language"`
	Alt      *string `json:"alt,omitempty"`
	Title    *string `json:"title,omitempty"`
}

type PromoSlideAdminItem struct {
	ID           uuid.UUID                    `json:"id"`
	ImagePath    string                       `json:"image_path"`
	ImageURL     string                       `json:"image_url"`
	OrderIndex   int                          `json:"order_index"`
	IsActive     bool                         `json:"is_active"`
	Translations []PromoSlideAdminTranslation `json:"translations"`
}
