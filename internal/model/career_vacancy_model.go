package model

import "github.com/google/uuid"

type CareerVacancy struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Employment  string    `json:"employment"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
}
