package admin

import (
	"context"
	"fmt"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
)

type ServiceMatrixUsecase interface {
	ListByService(ctx context.Context, service string, lang string) ([]model.ServiceMatrixAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceMatrixAdminItem, error)
	Create(ctx context.Context, input model.ServiceMatrixUpsertInput) (*model.ServiceMatrixAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.ServiceMatrixUpsertInput) (*model.ServiceMatrixAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type serviceMatrixUsecaseImpl struct {
	repo repository.ServiceMatrixRepository
}

func NewServiceMatrixUsecase(repo repository.ServiceMatrixRepository) ServiceMatrixUsecase {
	return &serviceMatrixUsecaseImpl{repo: repo}
}

func (u *serviceMatrixUsecaseImpl) ListByService(ctx context.Context, service string, lang string) ([]model.ServiceMatrixAdminItem, error) {
	items, err := u.repo.FindByService(ctx, strings.ToUpper(strings.TrimSpace(service)))
	if err != nil {
		return nil, err
	}
	res := make([]model.ServiceMatrixAdminItem, 0, len(items))
	for _, it := range items {
		res = append(res, toServiceMatrixAdminItem(it, lang))
	}
	return res, nil
}

func (u *serviceMatrixUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceMatrixAdminItem, error) {
	it, err := u.repo.FindByIDFull(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toServiceMatrixAdminItem(*it, lang)
	return &res, nil
}

func (u *serviceMatrixUsecaseImpl) Create(ctx context.Context, input model.ServiceMatrixUpsertInput) (*model.ServiceMatrixAdminItem, error) {
	service := strings.ToUpper(strings.TrimSpace(input.Service))
	lang := normServiceMatrixLang(input.Language)

	item := &entity.ServiceMatrix{
		Service:  service,
		Compact:  input.Compact,
		IsActive: input.IsActive,
	}
	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if err := u.upsertMatrixTranslations(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}

	if err := u.syncMatrixStructure(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}

	fresh, err := u.repo.FindByIDFull(ctx, item.ID)
	if err != nil {
		res := toServiceMatrixAdminItem(*item, lang)
		return &res, nil
	}
	res := toServiceMatrixAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceMatrixUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.ServiceMatrixUpsertInput) (*model.ServiceMatrixAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lang := normServiceMatrixLang(input.Language)
	if s := strings.ToUpper(strings.TrimSpace(input.Service)); s != "" {
		item.Service = s
	}
	item.Compact = input.Compact
	item.IsActive = input.IsActive
	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	if err := u.upsertMatrixTranslations(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}

	if err := u.syncMatrixStructure(ctx, item.ID, input, lang); err != nil {
		return nil, err
	}

	fresh, err := u.repo.FindByIDFull(ctx, item.ID)
	if err != nil {
		res := toServiceMatrixAdminItem(*item, lang)
		return &res, nil
	}
	res := toServiceMatrixAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceMatrixUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteByID(ctx, id)
}

func (u *serviceMatrixUsecaseImpl) upsertMatrixTranslations(ctx context.Context, matrixID uuid.UUID, input model.ServiceMatrixUpsertInput, lang string) error {
	if err := u.repo.UpsertTranslation(ctx, &entity.ServiceMatrixTranslation{
		MatrixID:    matrixID,
		Language:    lang,
		Title:       trimPtr(input.Title),
		Description: trimPtr(input.Description),
		Footnote:    trimPtr(input.Footnote),
	}); err != nil {
		return err
	}

	if lang == "ID" {
		enTitle := translateTextPtrServiceMatrix(input.Title, "id", "en")
		enDescription := translateTextPtrServiceMatrix(input.Description, "id", "en")
		enFootnote := translateTextPtrServiceMatrix(input.Footnote, "id", "en")
		_ = u.repo.UpsertTranslation(ctx, &entity.ServiceMatrixTranslation{
			MatrixID:    matrixID,
			Language:    "EN",
			Title:       enTitle,
			Description: enDescription,
			Footnote:    enFootnote,
		})
	}
	return nil
}

func (u *serviceMatrixUsecaseImpl) syncMatrixStructure(ctx context.Context, matrixID uuid.UUID, input model.ServiceMatrixUpsertInput, lang string) error {
	if err := u.repo.ClearStructure(ctx, matrixID); err != nil {
		return err
	}

	columnMap := map[string]uuid.UUID{}
	columns := make([]entity.ServiceMatrixColumn, 0, len(input.Columns))
	colTrLang := make([]entity.ServiceMatrixColumnTranslation, 0, len(input.Columns))
	colTrEN := make([]entity.ServiceMatrixColumnTranslation, 0, len(input.Columns))

	for i, c := range input.Columns {
		key := strings.TrimSpace(c.Key)
		if key == "" {
			key = fmt.Sprintf("col_%d", i+1)
		}
		id := uuid.New()
		label := strings.TrimSpace(c.Label)
		if label == "" {
			label = key
		}
		orderIndex := c.OrderIndex
		if orderIndex <= 0 {
			orderIndex = i + 1
		}

		columns = append(columns, entity.ServiceMatrixColumn{
			ID:         id,
			MatrixID:   matrixID,
			ColumnKey:  key,
			Popular:    c.Popular,
			OrderIndex: orderIndex,
		})
		colTrLang = append(colTrLang, entity.ServiceMatrixColumnTranslation{
			ColumnID: id,
			Language: lang,
			Label:    label,
		})
		if lang == "ID" {
			enLabel := translateTextValServiceMatrix(label, "id", "en")
			colTrEN = append(colTrEN, entity.ServiceMatrixColumnTranslation{ColumnID: id, Language: "EN", Label: enLabel})
		}
		columnMap[strings.ToLower(key)] = id
	}
	if err := u.repo.CreateColumns(ctx, columns); err != nil {
		return err
	}
	if err := u.repo.CreateColumnTranslations(ctx, colTrLang); err != nil {
		return err
	}
	if err := u.repo.CreateColumnTranslations(ctx, colTrEN); err != nil {
		return err
	}

	rows := make([]entity.ServiceMatrixRow, 0, len(input.Rows))
	rowTrLang := make([]entity.ServiceMatrixRowTranslation, 0, len(input.Rows))
	rowTrEN := make([]entity.ServiceMatrixRowTranslation, 0, len(input.Rows))
	cells := make([]entity.ServiceMatrixCell, 0)
	cellTrLang := make([]entity.ServiceMatrixCellTranslation, 0)
	cellTrEN := make([]entity.ServiceMatrixCellTranslation, 0)

	for i, r := range input.Rows {
		rowKey := strings.TrimSpace(r.Key)
		if rowKey == "" {
			rowKey = fmt.Sprintf("row_%d", i+1)
		}
		rowID := uuid.New()
		feature := strings.TrimSpace(r.Feature)
		if feature == "" {
			feature = rowKey
		}
		orderIndex := r.OrderIndex
		if orderIndex <= 0 {
			orderIndex = i + 1
		}

		rows = append(rows, entity.ServiceMatrixRow{ID: rowID, MatrixID: matrixID, RowKey: rowKey, OrderIndex: orderIndex})
		rowTrLang = append(rowTrLang, entity.ServiceMatrixRowTranslation{RowID: rowID, Language: lang, Feature: feature})
		if lang == "ID" {
			enFeature := translateTextValServiceMatrix(feature, "id", "en")
			rowTrEN = append(rowTrEN, entity.ServiceMatrixRowTranslation{RowID: rowID, Language: "EN", Feature: enFeature})
		}

		for _, cellIn := range r.Cells {
			colID, ok := columnMap[strings.ToLower(strings.TrimSpace(cellIn.ColumnKey))]
			if !ok {
				continue
			}
			cellID := uuid.New()
			cell := entity.ServiceMatrixCell{ID: cellID, RowID: rowID, ColumnID: colID, ValueBoolean: cellIn.ValueBoolean, ValueNumber: cellIn.ValueNumber}
			cells = append(cells, cell)
			if cellIn.ValueBoolean == nil && cellIn.ValueNumber == nil {
				text := strings.TrimSpace(ptrStr(cellIn.ValueText))
				if text == "" {
					text = "-"
				}
				cellTrLang = append(cellTrLang, entity.ServiceMatrixCellTranslation{CellID: cellID, Language: lang, ValueText: text})
				if lang == "ID" {
					enText := translateTextValServiceMatrix(text, "id", "en")
					cellTrEN = append(cellTrEN, entity.ServiceMatrixCellTranslation{CellID: cellID, Language: "EN", ValueText: enText})
				}
			}
		}
	}

	if err := u.repo.CreateRows(ctx, rows); err != nil {
		return err
	}
	if err := u.repo.CreateRowTranslations(ctx, rowTrLang); err != nil {
		return err
	}
	if err := u.repo.CreateRowTranslations(ctx, rowTrEN); err != nil {
		return err
	}
	if err := u.repo.CreateCells(ctx, cells); err != nil {
		return err
	}
	if err := u.repo.CreateCellTranslations(ctx, cellTrLang); err != nil {
		return err
	}
	if err := u.repo.CreateCellTranslations(ctx, cellTrEN); err != nil {
		return err
	}

	return nil
}

func normServiceMatrixLang(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" {
		return "ID"
	}
	return l
}

func toServiceMatrixAdminItem(it entity.ServiceMatrix, lang string) model.ServiceMatrixAdminItem {
	translations := make([]model.ServiceMatrixAdminTranslation, 0, len(it.Translations))
	for _, tr := range it.Translations {
		translations = append(translations, model.ServiceMatrixAdminTranslation{
			Language:    tr.Language,
			Title:       tr.Title,
			Description: tr.Description,
			Footnote:    tr.Footnote,
		})
	}
	var title, description, footnote *string
	for _, tr := range it.Translations {
		if strings.EqualFold(tr.Language, lang) {
			title, description, footnote = tr.Title, tr.Description, tr.Footnote
			break
		}
	}
	if title == nil && len(it.Translations) > 0 {
		title, description, footnote = it.Translations[0].Title, it.Translations[0].Description, it.Translations[0].Footnote
	}

	columns := make([]model.ServiceMatrixAdminColumnInput, 0, len(it.Columns))
	for _, col := range it.Columns {
		label := col.ColumnKey
		for _, tr := range col.Translations {
			if strings.EqualFold(tr.Language, lang) {
				label = tr.Label
				break
			}
		}
		columns = append(columns, model.ServiceMatrixAdminColumnInput{Key: col.ColumnKey, Label: label, Popular: col.Popular, OrderIndex: col.OrderIndex})
	}

	rows := make([]model.ServiceMatrixAdminRowInput, 0, len(it.Rows))
	for _, row := range it.Rows {
		feature := row.RowKey
		for _, tr := range row.Translations {
			if strings.EqualFold(tr.Language, lang) {
				feature = tr.Feature
				break
			}
		}
		cells := make([]model.ServiceMatrixAdminCellInput, 0, len(row.Cells))
		for _, cell := range row.Cells {
			text := ""
			for _, tr := range cell.Translations {
				if strings.EqualFold(tr.Language, lang) {
					text = tr.ValueText
					break
				}
			}
			var textPtr *string
			if text != "" {
				textPtr = &text
			}
			cells = append(cells, model.ServiceMatrixAdminCellInput{ColumnKey: findColumnKeyByID(it.Columns, cell.ColumnID), ValueBoolean: cell.ValueBoolean, ValueNumber: cell.ValueNumber, ValueText: textPtr})
		}
		rows = append(rows, model.ServiceMatrixAdminRowInput{Key: row.RowKey, Feature: feature, OrderIndex: row.OrderIndex, Cells: cells})
	}

	return model.ServiceMatrixAdminItem{
		ID:           it.ID,
		Service:      it.Service,
		Compact:      it.Compact,
		IsActive:     it.IsActive,
		Title:        title,
		Description:  description,
		Footnote:     footnote,
		Columns:      columns,
		Rows:         rows,
		Translations: translations,
	}
}

func findColumnKeyByID(columns []entity.ServiceMatrixColumn, id uuid.UUID) string {
	for _, c := range columns {
		if c.ID == id {
			return c.ColumnKey
		}
	}
	return ""
}

func translateTextPtrServiceMatrix(text *string, fromLang, toLang string) *string {
	if text == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*text)
	if trimmed == "" {
		return nil
	}
	translated, err := gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{From: fromLang, To: toLang})
	if err != nil || strings.TrimSpace(translated) == "" {
		fallback := trimmed
		return &fallback
	}
	res := strings.TrimSpace(translated)
	return &res
}

func translateTextValServiceMatrix(text, fromLang, toLang string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}
	translated, err := gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{From: fromLang, To: toLang})
	if err != nil || strings.TrimSpace(translated) == "" {
		return trimmed
	}
	return strings.TrimSpace(translated)
}

func ptrStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
