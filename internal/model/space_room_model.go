package model

import "github.com/google/uuid"

type SpaceRoomCard struct {
	ID uuid.UUID `json:"id"`

	Slug string `json:"slug"`

	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`

	MainImageURL   *string `json:"main_image_url,omitempty"`
	MainImageAlt   *string `json:"main_image_alt,omitempty"`
	MainImageTitle *string `json:"main_image_title,omitempty"`

	Capacity *int `json:"capacity,omitempty"`
	Floor    *int `json:"floor,omitempty"`

	Facilities []string `json:"facilities"`

	Rating     *float64 `json:"rating,omitempty"`
	RatingText *string  `json:"rating_text,omitempty"`

	IsAvailable bool   `json:"is_available"`
	StatusText  string `json:"status_text"`

	Price *float64 `json:"price,omitempty"`

	Tags []string `json:"tags"`

	ActionLabel string `json:"action_label"`
	ActionState string `json:"action_state"`
}
