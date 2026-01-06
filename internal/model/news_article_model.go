package model

import "github.com/google/uuid"

type NewsArticle struct {
	ID uuid.UUID `json:"id"`

	ImageURL *string `json:"image_url,omitempty"`

	IsActive    bool   `json:"is_active"`
	PublishedAt string `json:"published_at"`

	ViewCount int `json:"view_count"`
	LikeCount int `json:"like_count"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`

	Meta NewsArticleMeta `json:"meta"`

	Categories []NewsArticleCategory `json:"categories"`
}

type NewsArticleMeta struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

type NewsArticleCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
}

type NewsCard struct {
	ID          uuid.UUID `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	ImageURL    *string   `json:"image_url,omitempty"`
	PublishedAt string    `json:"published_at"`
}
