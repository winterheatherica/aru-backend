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

type NewsCategoryUpsertInput struct {
	Language    string  `json:"language"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	IsActive    bool    `json:"is_active"`
}

type NewsCategoryAdminTranslation struct {
	Language    string  `json:"language"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
}

type NewsCategoryAdminItem struct {
	ID           uuid.UUID                      `json:"id"`
	Name         *string                        `json:"name,omitempty"`
	Slug         *string                        `json:"slug,omitempty"`
	Description  *string                        `json:"description,omitempty"`
	IsActive     bool                           `json:"is_active"`
	Translations []NewsCategoryAdminTranslation `json:"translations"`
}
