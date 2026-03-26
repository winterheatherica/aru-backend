package model

import "github.com/google/uuid"

type CareerVacancy struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Employment  string    `json:"employment"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}

type CareerVacancyUpsertInput struct {
	Language    string  `json:"language"`
	Title       string  `json:"title"`
	Employment  string  `json:"employment"`
	Location    string  `json:"location"`
	Description *string `json:"description"`
	OpenedAt    string  `json:"opened_at"`
	ClosedAt    *string `json:"closed_at"`
	IsActive    bool    `json:"is_active"`
}

type CareerVacancyAdminTranslation struct {
	Language    string  `json:"language"`
	Description *string `json:"description,omitempty"`
}

type CareerVacancyAdminItem struct {
	ID           uuid.UUID                       `json:"id"`
	Title        string                          `json:"title"`
	Slug         string                          `json:"slug"`
	Employment   string                          `json:"employment"`
	Location     string                          `json:"location"`
	IsActive     bool                            `json:"is_active"`
	OpenedAt     string                          `json:"opened_at"`
	ClosedAt     *string                         `json:"closed_at,omitempty"`
	Translations []CareerVacancyAdminTranslation `json:"translations"`
}
