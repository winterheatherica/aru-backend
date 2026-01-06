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
