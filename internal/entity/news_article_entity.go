package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type NewsArticle struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ImagePath *string `gorm:"column:image_path" json:"image_path,omitempty"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	PublishedAt time.Time `gorm:"column:published_at" json:"published_at"`

	ViewCount int `gorm:"column:view_count;default:0" json:"view_count"`
	LikeCount int `gorm:"column:like_count;default:0" json:"like_count"`

	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []NewsArticleTranslation `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"translations"`

	Categories []NewsCategory `gorm:"many2many:news_article_categories;joinForeignKey:ArticleID;joinReferences:CategoryID" json:"categories,omitempty"`
}

func (NewsArticle) TableName() string {
	return "news_articles"
}

type NewsArticleTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ArticleID uuid.UUID `gorm:"column:article_id;type:uuid;not null" json:"article_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Slug    string `gorm:"column:slug;not null" json:"slug"`
	Title   string `gorm:"column:title;not null" json:"title"`
	Content string `gorm:"column:content;not null" json:"content"`

	MetaTitle       *string        `gorm:"column:meta_title" json:"meta_title,omitempty"`
	MetaDescription *string        `gorm:"column:meta_description" json:"meta_description,omitempty"`
	MetaKeywords    pq.StringArray `gorm:"column:meta_keywords;type:text[]" json:"meta_keywords,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (NewsArticleTranslation) TableName() string {
	return "news_article_translations"
}
