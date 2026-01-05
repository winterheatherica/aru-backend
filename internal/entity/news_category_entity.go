package entity

import (
	"time"

	"github.com/google/uuid"
)

type NewsCategory struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []NewsCategoryTranslation `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (NewsCategory) TableName() string {
	return "news_categories"
}

type NewsCategoryTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CategoryID uuid.UUID `gorm:"column:category_id;type:uuid;not null" json:"category_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Name        string  `gorm:"column:name;not null" json:"name"`
	Slug        string  `gorm:"column:slug;not null" json:"slug"`
	Description *string `gorm:"column:description" json:"description,omitempty"`

	MetaTitle       *string  `gorm:"column:meta_title" json:"meta_title,omitempty"`
	MetaDescription *string  `gorm:"column:meta_description" json:"meta_description,omitempty"`
	MetaKeywords    []string `gorm:"column:meta_keywords;type:text[]" json:"meta_keywords,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (NewsCategoryTranslation) TableName() string {
	return "news_category_translations"
}
