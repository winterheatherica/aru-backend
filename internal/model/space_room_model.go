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

	IsAvailable bool   `json:"is_available"`
	StatusText  string `json:"status_text"`

	Price *float64 `json:"price,omitempty"`

	ActionLabel string `json:"action_label"`
	ActionState string `json:"action_state"`
}

type SpaceRoomDetail struct {
	ID uuid.UUID `json:"id"`

	Slug string `json:"slug"`

	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`

	Images []SpaceRoomImage `json:"images"`

	Capacity *int `json:"capacity,omitempty"`
	Floor    *int `json:"floor,omitempty"`

	Facilities []string `json:"facilities"`

	IsAvailable bool   `json:"is_available"`
	StatusText  string `json:"status_text"`

	Price *float64 `json:"price,omitempty"`

	ActionLabel string `json:"action_label"`
	ActionState string `json:"action_state"`
}

type SpaceRoomImage struct {
	URL   string  `json:"url"`
	Alt   *string `json:"alt,omitempty"`
	Title *string `json:"title,omitempty"`
}
