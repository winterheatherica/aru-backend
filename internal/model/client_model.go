package model

import "github.com/google/uuid"

type Client struct {
	ID    uuid.UUID `json:"id"`
	Src   string    `json:"src"`
	Alt   string    `json:"alt"`
	Title *string   `json:"title,omitempty"`
	Desc  *string   `json:"description,omitempty"`
	Order int       `json:"order"`
}

type ClientUpsertInput struct {
	Language               string  `json:"language"`
	Alt                    *string `json:"alt"`
	Title                  *string `json:"title"`
	Description            *string `json:"description"`
	OrderIndex             int     `json:"order_index"`
	IsActiveClientScroller bool    `json:"is_active_client_scroller"`
}

type ClientAdminTranslation struct {
	Language    string  `json:"language"`
	Alt         *string `json:"alt,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

type ClientAdminItem struct {
	ID                     uuid.UUID                `json:"id"`
	ImagePath              string                   `json:"image_path"`
	ImageURL               string                   `json:"image_url"`
	OrderIndex             int                      `json:"order_index"`
	IsActiveClientScroller bool                     `json:"is_active_client_scroller"`
	Translations           []ClientAdminTranslation `json:"translations"`
}
