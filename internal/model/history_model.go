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
