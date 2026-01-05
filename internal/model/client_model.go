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
