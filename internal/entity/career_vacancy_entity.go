package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CareerVacancy struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Title      string `gorm:"column:title;not null" json:"title"`
	Slug       string `gorm:"column:slug;not null;unique" json:"slug"`
	Employment string `gorm:"column:employment;type:employment;not null" json:"employment"`
	Location   string `gorm:"column:location;not null" json:"location"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	OpenedAt time.Time  `gorm:"column:opened_at" json:"opened_at"`
	ClosedAt *time.Time `gorm:"column:closed_at" json:"closed_at,omitempty"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []CareerVacancyTranslation `gorm:"foreignKey:VacancyID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (CareerVacancy) TableName() string {
	return "career_vacancies"
}

type CareerVacancyTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	VacancyID uuid.UUID `gorm:"column:vacancy_id;type:uuid;not null" json:"vacancy_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Description     *string        `gorm:"column:description" json:"description,omitempty"`
	MetaTitle       *string        `gorm:"column:meta_title" json:"meta_title,omitempty"`
	MetaDescription *string        `gorm:"column:meta_description" json:"meta_description,omitempty"`
	MetaKeywords    pq.StringArray `gorm:"column:meta_keywords;type:text[]" json:"meta_keywords,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (CareerVacancyTranslation) TableName() string {
	return "career_vacancy_translations"
}
