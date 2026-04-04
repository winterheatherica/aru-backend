package model

import (
	"time"

	"github.com/google/uuid"
)

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

	Price *float64 `json:"price,omitempty"`
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

	Price *float64 `json:"price,omitempty"`
}

type SpaceRoomImage struct {
	URL   string  `json:"url"`
	Alt   *string `json:"alt,omitempty"`
	Title *string `json:"title,omitempty"`
}

type SpaceRoomAdminImageInput struct {
	Alt   string `json:"alt"`
	Title string `json:"title"`
}

type SpaceRoomAdminUpsertInput struct {
	Language    string     `json:"language"`
	Code        string     `json:"code"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Facilities  []string   `json:"facilities"`
	MetaKeywords []string  `json:"meta_keywords"`
	Capacity    *int       `json:"capacity"`
	Floor       *int       `json:"floor"`
	Address     *string    `json:"address"`
	IsActive    bool       `json:"is_active"`
	UploadedBy  *uuid.UUID `json:"uploaded_by,omitempty"`
	Images      []SpaceRoomAdminImageInput `json:"images"`
	ThumbnailIndex int `json:"thumbnail_index"`
}

type SpaceRoomAdminTranslation struct {
	Language        string   `json:"language"`
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Description     *string  `json:"description,omitempty"`
	Facilities      []string `json:"facilities,omitempty"`
	MetaTitle       *string  `json:"meta_title,omitempty"`
	MetaDescription *string  `json:"meta_description,omitempty"`
	MetaKeywords    []string `json:"meta_keywords,omitempty"`
}

type SpaceRoomAdminImage struct {
	ID         uuid.UUID `json:"id"`
	ImageURL   string    `json:"image_url"`
	OrderIndex int       `json:"order_index"`
	IsActive   bool      `json:"is_active"`
	IsThumbnail bool     `json:"is_thumbnail"`
	Alt        *string   `json:"alt,omitempty"`
	Title      *string   `json:"title,omitempty"`
}

type SpaceRoomAdminItem struct {
	ID             uuid.UUID `json:"id"`
	Code           string    `json:"code"`
	Capacity       *int      `json:"capacity,omitempty"`
	Floor          *int      `json:"floor,omitempty"`
	Address        *string   `json:"address,omitempty"`
	MainImageURL   *string   `json:"main_image_url,omitempty"`
	UploadedBy     *uuid.UUID `json:"uploaded_by,omitempty"`
	UploadedByName *string    `json:"uploaded_by_name,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Translations   []SpaceRoomAdminTranslation `json:"translations"`
	Images         []SpaceRoomAdminImage `json:"images"`
}
