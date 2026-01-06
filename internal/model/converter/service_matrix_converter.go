package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func ServiceMatrixToModel(
	matrix entity.ServiceMatrix,
	lang string,
) *model.ServiceMatrix {

	mt := findServiceMatrixTranslation(matrix, lang)

	result := &model.ServiceMatrix{
		ID:      matrix.ID,
		Service: matrix.Service,
		Compact: matrix.Compact,

		Title:       nil,
		Description: nil,
		Footnote:    nil,

		Columns: make([]model.ServiceMatrixColumn, 0),
		Rows:    make([]model.ServiceMatrixRow, 0),
	}

	if mt != nil {
		result.Title = mt.Title
		result.Description = mt.Description
		result.Footnote = mt.Footnote
	}

	for _, c := range matrix.Columns {
		col := serviceMatrixColumnToModel(c, lang)
		if col != nil {
			result.Columns = append(result.Columns, *col)
		}
	}

	for _, r := range matrix.Rows {
		row := serviceMatrixRowToModel(r, lang)
		if row != nil {
			result.Rows = append(result.Rows, *row)
		}
	}

	return result
}

func serviceMatrixColumnToModel(
	column entity.ServiceMatrixColumn,
	lang string,
) *model.ServiceMatrixColumn {

	ct := findServiceMatrixColumnTranslation(column, lang)
	if ct == nil {
		return nil
	}

	return &model.ServiceMatrixColumn{
		ID: column.ID,

		Key:   column.ColumnKey,
		Label: ct.Label,

		Popular:    column.Popular,
		OrderIndex: column.OrderIndex,
	}
}

func serviceMatrixRowToModel(
	row entity.ServiceMatrixRow,
	lang string,
) *model.ServiceMatrixRow {

	rt := findServiceMatrixRowTranslation(row, lang)
	if rt == nil {
		return nil
	}

	mRow := &model.ServiceMatrixRow{
		ID: row.ID,

		Key: row.RowKey,

		Feature: rt.Feature,

		OrderIndex: row.OrderIndex,

		Cells: make([]model.ServiceMatrixCell, 0),
	}

	for _, c := range row.Cells {
		cell := serviceMatrixCellToModel(c, lang)
		if cell != nil {
			mRow.Cells = append(mRow.Cells, *cell)
		}
	}

	return mRow
}

func serviceMatrixCellToModel(
	cell entity.ServiceMatrixCell,
	lang string,
) *model.ServiceMatrixCell {

	mCell := &model.ServiceMatrixCell{
		RowID:    cell.RowID,
		ColumnID: cell.ColumnID,

		ValueBoolean: nil,
		ValueNumber:  nil,
		ValueText:    nil,
	}

	if cell.ValueBoolean != nil {
		mCell.ValueBoolean = cell.ValueBoolean
		return mCell
	}

	if cell.ValueNumber != nil {
		mCell.ValueNumber = cell.ValueNumber
		return mCell
	}

	ct := findServiceMatrixCellTranslation(cell, lang)
	if ct != nil {
		mCell.ValueText = &ct.ValueText
	}

	return mCell
}

func findServiceMatrixTranslation(
	matrix entity.ServiceMatrix,
	lang string,
) *entity.ServiceMatrixTranslation {

	for i := range matrix.Translations {
		if matrix.Translations[i].Language == lang {
			return &matrix.Translations[i]
		}
	}
	return nil
}

func findServiceMatrixColumnTranslation(
	column entity.ServiceMatrixColumn,
	lang string,
) *entity.ServiceMatrixColumnTranslation {

	for i := range column.Translations {
		if column.Translations[i].Language == lang {
			return &column.Translations[i]
		}
	}
	return nil
}

func findServiceMatrixRowTranslation(
	row entity.ServiceMatrixRow,
	lang string,
) *entity.ServiceMatrixRowTranslation {

	for i := range row.Translations {
		if row.Translations[i].Language == lang {
			return &row.Translations[i]
		}
	}
	return nil
}

func findServiceMatrixCellTranslation(
	cell entity.ServiceMatrixCell,
	lang string,
) *entity.ServiceMatrixCellTranslation {

	for i := range cell.Translations {
		if cell.Translations[i].Language == lang {
			return &cell.Translations[i]
		}
	}
	return nil
}
