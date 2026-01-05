package model

import "github.com/google/uuid"

type NewsCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`

	Meta NewsCategoryMeta `json:"meta"`
}

type NewsCategoryMeta struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}
