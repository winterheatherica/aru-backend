package model

import "github.com/google/uuid"

type ServiceMatrix struct {
	ID uuid.UUID `json:"id"`

	Service string `json:"service"`
	Compact bool   `json:"compact"`

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Footnote    *string `json:"footnote,omitempty"`

	Columns []ServiceMatrixColumn `json:"columns"`
	Rows    []ServiceMatrixRow    `json:"rows"`
}

type ServiceMatrixColumn struct {
	ID uuid.UUID `json:"id"`

	Key string `json:"key"`

	Label string `json:"label"`

	Popular    bool `json:"popular"`
	OrderIndex int  `json:"order_index"`
}

type ServiceMatrixRow struct {
	ID uuid.UUID `json:"id"`

	Key string `json:"key"`

	Feature string `json:"feature"`

	OrderIndex int `json:"order_index"`

	Cells []ServiceMatrixCell `json:"cells"`
}

type ServiceMatrixCell struct {
	RowID    uuid.UUID `json:"row_id"`
	ColumnID uuid.UUID `json:"column_id"`

	ValueBoolean *bool    `json:"value_boolean,omitempty"`
	ValueNumber  *float64 `json:"value_number,omitempty"`
	ValueText    *string  `json:"value_text,omitempty"`
}

type ServiceMatrixAdminTranslation struct {
	Language    string  `json:"language"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Footnote    *string `json:"footnote,omitempty"`
}

type ServiceMatrixAdminColumnInput struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	Popular    bool   `json:"popular"`
	OrderIndex int    `json:"order_index"`
}

type ServiceMatrixAdminCellInput struct {
	ColumnKey    string   `json:"column_key"`
	ValueBoolean *bool    `json:"value_boolean"`
	ValueNumber  *float64 `json:"value_number"`
	ValueText    *string  `json:"value_text"`
}

type ServiceMatrixAdminRowInput struct {
	Key        string                        `json:"key"`
	Feature    string                        `json:"feature"`
	OrderIndex int                           `json:"order_index"`
	Cells      []ServiceMatrixAdminCellInput `json:"cells"`
}

type ServiceMatrixAdminItem struct {
	ID uuid.UUID `json:"id"`

	Service  string `json:"service"`
	Compact  bool   `json:"compact"`
	IsActive bool   `json:"is_active"`

	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Footnote    *string `json:"footnote,omitempty"`

	Columns      []ServiceMatrixAdminColumnInput `json:"columns"`
	Rows         []ServiceMatrixAdminRowInput    `json:"rows"`
	Translations []ServiceMatrixAdminTranslation `json:"translations"`
}

type ServiceMatrixUpsertInput struct {
	Service     string  `json:"service"`
	Language    string  `json:"language"`
	Compact     bool    `json:"compact"`
	IsActive    bool    `json:"is_active"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Footnote    *string `json:"footnote"`

	Columns []ServiceMatrixAdminColumnInput `json:"columns"`
	Rows    []ServiceMatrixAdminRowInput    `json:"rows"`
}
