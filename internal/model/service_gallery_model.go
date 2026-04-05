package model

import "github.com/google/uuid"

type ServiceGallery struct {
	ID uuid.UUID `json:"id"`

	Service   string `json:"service"`
	MediaType string `json:"media_type"`

	Src       string  `json:"src"`
	Thumbnail *string `json:"thumbnail,omitempty"`

	Title   *string `json:"title,omitempty"`
	Alt     *string `json:"alt,omitempty"`
	Caption *string `json:"caption,omitempty"`
}

type ServiceGalleryUpsertInput struct {
	Service   string  `json:"service"`
	MediaType string  `json:"media_type"`
	Language  string  `json:"language"`
	Title     *string `json:"title"`
	Alt       *string `json:"alt"`
	Caption   *string `json:"caption"`
	IsActive  bool    `json:"is_active"`
}

type ServiceGalleryAdminTranslation struct {
	Language string  `json:"language"`
	Title    *string `json:"title,omitempty"`
	Alt      *string `json:"alt,omitempty"`
	Caption  *string `json:"caption,omitempty"`
}

type ServiceGalleryAdminItem struct {
	ID           uuid.UUID                        `json:"id"`
	Service      string                           `json:"service"`
	MediaType    string                           `json:"media_type"`
	Src          string                           `json:"src"`
	ImageURL     string                           `json:"image_url"`
	Thumbnail    *string                          `json:"thumbnail,omitempty"`
	IsActive     bool                             `json:"is_active"`
	Title        *string                          `json:"title,omitempty"`
	Alt          *string                          `json:"alt,omitempty"`
	Caption      *string                          `json:"caption,omitempty"`
	Translations []ServiceGalleryAdminTranslation `json:"translations"`
}
