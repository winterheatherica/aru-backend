package model

import "github.com/google/uuid"

type History struct {
	ID          uuid.UUID `json:"id"`
	Year        *int      `json:"year,omitempty"`
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`

	TableHeaders []string   `json:"table_headers,omitempty"`
	TableRows    [][]string `json:"table_rows,omitempty"`
}

type HistoryUpsertInput struct {
	Language     string     `json:"language"`
	Year         *int       `json:"year"`
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	TableHeaders []string   `json:"table_headers"`
	TableRows    [][]string `json:"table_rows"`
	IsActive     bool       `json:"is_active"`
}

type HistoryAdminItem struct {
	ID           uuid.UUID  `json:"id"`
	Language     string     `json:"language"`
	Year         *int       `json:"year,omitempty"`
	Title        *string    `json:"title,omitempty"`
	Description  *string    `json:"description,omitempty"`
	TableHeaders []string   `json:"table_headers,omitempty"`
	TableRows    [][]string `json:"table_rows,omitempty"`
	IsActive     bool       `json:"is_active"`
}
