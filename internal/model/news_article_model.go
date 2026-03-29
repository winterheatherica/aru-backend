package model

import (
	"time"

	"github.com/google/uuid"
)

type NewsArticle struct {
	ID uuid.UUID `json:"id"`

	ImageURL *string `json:"image_url,omitempty"`

	IsActive    bool   `json:"is_active"`
	PublishedAt         string  `json:"published_at"`
	PublishedBy         string  `json:"published_by,omitempty"`
	PublishedByAvatarURL *string `json:"published_by_avatar_url,omitempty"`

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
	ImagePath   *string   `json:"image_path,omitempty"`
	PublishedAt string    `json:"published_at"`
}

type NewsArticleUpsertInput struct {
	Language        string     `json:"language"`
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	MetaTitle       *string    `json:"meta_title"`
	MetaDescription *string    `json:"meta_description"`
	MetaKeywords    []string   `json:"meta_keywords"`
	CategoryIDs     []string   `json:"category_ids"`
	UploadedBy      *uuid.UUID `json:"uploaded_by,omitempty"`
	PublishedAt     time.Time  `json:"published_at"`
	IsActive        bool       `json:"is_active"`
}

type NewsArticleAdminTranslation struct {
	Language        string   `json:"language"`
	Slug            string   `json:"slug"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	MetaTitle       *string  `json:"meta_title,omitempty"`
	MetaDescription *string  `json:"meta_description,omitempty"`
	MetaKeywords    []string `json:"meta_keywords,omitempty"`
}

type NewsArticleAdminCategory struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type NewsArticleAdminItem struct {
	ID             uuid.UUID                     `json:"id"`
	ImagePath      *string                       `json:"image_path,omitempty"`
	ImageURL       *string                       `json:"image_url,omitempty"`
	IsActive       bool                          `json:"is_active"`
	PublishedAt    time.Time                     `json:"published_at"`
	ViewCount      int                           `json:"view_count"`
	LikeCount      int                           `json:"like_count"`
	UploadedBy     *uuid.UUID                    `json:"uploaded_by,omitempty"`
	UploadedByName *string                       `json:"uploaded_by_name,omitempty"`
	CategoryIDs    []uuid.UUID                   `json:"category_ids,omitempty"`
	Categories     []NewsArticleAdminCategory    `json:"categories,omitempty"`
	Translations   []NewsArticleAdminTranslation `json:"translations"`
}
