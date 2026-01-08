package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
)

type History struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Year *int `gorm:"column:year" json:"year,omitempty"`

	Title       *string `gorm:"column:title" json:"title,omitempty"`
	Description *string `gorm:"column:description" json:"description,omitempty"`

	TableHeaders pq.StringArray   `gorm:"column:table_headers;type:text[]" json:"table_headers,omitempty"`
	TableRows    pgtype.TextArray `gorm:"column:table_rows;type:text[][]" json:"-"`

	IsActive bool `gorm:"column:is_active;default:true" json:"is_active"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (History) TableName() string {
	return "histories"
}
