package entity

import (
	"time"

	"github.com/google/uuid"
)

type ServiceMatrix struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	Service string `gorm:"column:service;type:service;not null" json:"service"`
	Compact bool   `gorm:"column:compact;default:false" json:"compact"`

	UploadedBy *uuid.UUID `gorm:"column:uploaded_by" json:"uploaded_by"`
	IsActive   bool       `gorm:"column:is_active;default:true" json:"is_active"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceMatrixTranslation `gorm:"foreignKey:MatrixID;constraint:OnDelete:CASCADE" json:"translations"`

	Columns []ServiceMatrixColumn `gorm:"foreignKey:MatrixID;constraint:OnDelete:CASCADE" json:"columns"`
	Rows    []ServiceMatrixRow    `gorm:"foreignKey:MatrixID;constraint:OnDelete:CASCADE" json:"rows"`
}

func (ServiceMatrix) TableName() string { return "service_matrices" }

type ServiceMatrixTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	MatrixID uuid.UUID `gorm:"column:matrix_id;type:uuid;not null" json:"matrix_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Title       *string `gorm:"column:title" json:"title,omitempty"`
	Description *string `gorm:"column:description" json:"description,omitempty"`
	Footnote    *string `gorm:"column:footnote" json:"footnote,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceMatrixTranslation) TableName() string {
	return "service_matrix_translations"
}

type ServiceMatrixColumn struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	MatrixID uuid.UUID `gorm:"column:matrix_id;type:uuid;not null" json:"matrix_id"`

	ColumnKey string `gorm:"column:column_key;not null" json:"column_key"`

	Popular    bool `gorm:"column:popular;default:false" json:"popular"`
	OrderIndex int  `gorm:"column:order_index;default:1" json:"order_index"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceMatrixColumnTranslation `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE" json:"translations"`
	Cells        []ServiceMatrixCell              `gorm:"foreignKey:ColumnID;constraint:OnDelete:CASCADE" json:"cells"`
}

func (ServiceMatrixColumn) TableName() string { return "service_matrix_columns" }

type ServiceMatrixColumnTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	ColumnID uuid.UUID `gorm:"column:column_id;type:uuid;not null" json:"column_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Label string `gorm:"column:label;not null" json:"label"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceMatrixColumnTranslation) TableName() string {
	return "service_matrix_column_translations"
}

type ServiceMatrixRow struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	MatrixID uuid.UUID `gorm:"column:matrix_id;type:uuid;not null" json:"matrix_id"`

	RowKey     string `gorm:"column:row_key;not null" json:"row_key"`
	OrderIndex int    `gorm:"column:order_index;default:1" json:"order_index"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceMatrixRowTranslation `gorm:"foreignKey:RowID;constraint:OnDelete:CASCADE" json:"translations"`
	Cells        []ServiceMatrixCell           `gorm:"foreignKey:RowID;constraint:OnDelete:CASCADE" json:"cells"`
}

func (ServiceMatrixRow) TableName() string { return "service_matrix_rows" }

type ServiceMatrixRowTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	RowID uuid.UUID `gorm:"column:row_id;type:uuid;not null" json:"row_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	Feature string `gorm:"column:feature;not null" json:"feature"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceMatrixRowTranslation) TableName() string {
	return "service_matrix_row_translations"
}

type ServiceMatrixCell struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	RowID    uuid.UUID `gorm:"column:row_id;type:uuid;not null" json:"row_id"`
	ColumnID uuid.UUID `gorm:"column:column_id;type:uuid;not null" json:"column_id"`

	ValueBoolean *bool    `gorm:"column:value_boolean" json:"value_boolean,omitempty"`
	ValueNumber  *float64 `gorm:"column:value_number" json:"value_number,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Translations []ServiceMatrixCellTranslation `gorm:"foreignKey:CellID;constraint:OnDelete:CASCADE" json:"translations"`
}

func (ServiceMatrixCell) TableName() string { return "service_matrix_cells" }

type ServiceMatrixCellTranslation struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`

	CellID uuid.UUID `gorm:"column:cell_id;type:uuid;not null" json:"cell_id"`

	Language string `gorm:"column:language;type:language;not null" json:"language"`

	ValueText string `gorm:"column:value_text;not null" json:"value_text"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ServiceMatrixCellTranslation) TableName() string {
	return "service_matrix_cell_translations"
}
