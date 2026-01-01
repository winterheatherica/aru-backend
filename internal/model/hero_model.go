package model

import "github.com/google/uuid"

type HeroSlide struct {
	ID       uuid.UUID `json:"id"`
	Src      string    `json:"src"`
	Alt      string    `json:"alt"`
	Title    *string   `json:"title,omitempty"`
	CtaLabel *string   `json:"ctaLabel,omitempty"`
	Banner   string    `json:"banner"`
	Order    int       `json:"order"`
}
