package admin

import (
	"context"
	"fmt"
	"log"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type HistoryUsecase interface {
	List(ctx context.Context) ([]model.HistoryAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.HistoryAdminItem, error)
	Create(ctx context.Context, input model.HistoryUpsertInput) (*model.HistoryAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.HistoryUpsertInput) (*model.HistoryAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type historyUsecaseImpl struct{ repo repository.HistoryRepository }

func NewHistoryUsecase(repo repository.HistoryRepository) HistoryUsecase {
	return &historyUsecaseImpl{repo: repo}
}

func normLang(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" {
		l = "ID"
	}
	return l
}

func (u *historyUsecaseImpl) List(ctx context.Context) ([]model.HistoryAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]model.HistoryAdminItem, 0, len(items))
	for _, it := range items {
		res = append(res, toHistoryAdminItem(it))
	}
	return res, nil
}

func (u *historyUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.HistoryAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := toHistoryAdminItem(*it)
	return &res, nil
}

func (u *historyUsecaseImpl) Create(ctx context.Context, input model.HistoryUpsertInput) (*model.HistoryAdminItem, error) {
	if input.Year == nil {
		return nil, fmt.Errorf("year is required")
	}
	lang := normLang(input.Language)

	if existing, err := u.repo.FindByYearAndLanguage(ctx, *input.Year, lang); err == nil && existing != nil {
		return nil, fmt.Errorf("history for year %d and language %s already exists", *input.Year, lang)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	headers, rows := normalizeHistoryMatrix(input.TableHeaders, input.TableRows)
	item := &entity.History{
		Language:     lang,
		Year:         input.Year,
		Title:        input.Title,
		Description:  input.Description,
		TableHeaders: headers,
		TableRows:    toPgText2D(rows),
		IsActive:     input.IsActive,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	if lang == "ID" {
		if _, err := u.repo.FindByYearAndLanguage(ctx, *input.Year, "EN"); err == gorm.ErrRecordNotFound {
			enTitle, fbTitle := translateTextPtrHistory(input.Title, "id", "en")
			enDesc, fbDesc := translateTextPtrHistory(input.Description, "id", "en")
			enHeaders, enRows, fbMatrix := translateMatrix(headers, rows, "id", "en")
			enItem := &entity.History{
				Language:          "EN",
				Year:              input.Year,
				Title:             enTitle,
				Description:       enDesc,
				TableHeaders:      enHeaders,
				TableRows:         toPgText2D(enRows),
				IsActive:          input.IsActive,
				IsMachineFallback: fbTitle || fbDesc || fbMatrix,
			}
			if cerr := u.repo.Create(ctx, enItem); cerr != nil {
				log.Printf("history: failed to create EN auto-translation for year %d: %v", *input.Year, cerr)
			}
		}
	}

	created, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := toHistoryAdminItem(*item)
		return &res, nil
	}
	res := toHistoryAdminItem(*created)
	return &res, nil
}

func (u *historyUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.HistoryUpsertInput) (*model.HistoryAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if input.Year == nil {
		return nil, fmt.Errorf("year is required")
	}
	lang := normLang(input.Language)

	if conflict, err := u.repo.FindByYearAndLanguage(ctx, *input.Year, lang); err == nil && conflict != nil && conflict.ID != id {
		return nil, fmt.Errorf("history for year %d and language %s already exists", *input.Year, lang)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	headers, rows := normalizeHistoryMatrix(input.TableHeaders, input.TableRows)
	item.Language = lang
	item.Year = input.Year
	item.Title = input.Title
	item.Description = input.Description
	item.TableHeaders = headers
	item.TableRows = toPgText2D(rows)
	item.IsActive = input.IsActive
	item.IsMachineFallback = false

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	updated, err := u.repo.FindByID(ctx, id)
	if err != nil {
		res := toHistoryAdminItem(*item)
		return &res, nil
	}
	res := toHistoryAdminItem(*updated)
	return &res, nil
}

func (u *historyUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteByID(ctx, id)
}

func toHistoryAdminItem(h entity.History) model.HistoryAdminItem {
	return model.HistoryAdminItem{
		ID:                h.ID,
		Language:          h.Language,
		Year:              h.Year,
		Title:             h.Title,
		Description:       h.Description,
		TableHeaders:      []string(h.TableHeaders),
		TableRows:         parsePgText2D(h.TableRows),
		IsActive:          h.IsActive,
		IsMachineFallback: h.IsMachineFallback,
	}
}

func parsePgText2D(arr pgtype.TextArray) [][]string {
	if arr.Status != pgtype.Present || len(arr.Dimensions) < 2 {
		return nil
	}
	cols := int(arr.Dimensions[1].Length)
	if cols <= 0 {
		return nil
	}
	rows := make([][]string, 0)
	for i := 0; i < len(arr.Elements); i += cols {
		row := make([]string, 0, cols)
		for j := 0; j < cols && i+j < len(arr.Elements); j++ {
			row = append(row, arr.Elements[i+j].String)
		}
		rows = append(rows, row)
	}
	return rows
}

func normalizeHistoryMatrix(headers []string, rows [][]string) ([]string, [][]string) {
	nHeaders := make([]string, 0, len(headers))
	for _, h := range headers {
		h = strings.TrimSpace(h)
		if h != "" {
			nHeaders = append(nHeaders, h)
		}
	}

	cols := len(nHeaders)
	if cols == 0 {
		for _, r := range rows {
			if len(r) > cols {
				cols = len(r)
			}
		}
	}
	if cols == 0 {
		return nHeaders, nil
	}

	// rescue horizontal flattened row
	if len(rows) == 1 && len(rows[0]) > cols && len(rows[0])%cols == 0 {
		flat := rows[0]
		chunked := make([][]string, 0, len(flat)/cols)
		for i := 0; i < len(flat); i += cols {
			chunked = append(chunked, append([]string{}, flat[i:i+cols]...))
		}
		rows = chunked
	}

	// rescue vertical flattened rows (Nx1)
	if cols > 1 && len(rows) >= cols {
		allSingle := true
		flat := make([]string, 0, len(rows))
		for _, r := range rows {
			if len(r) != 1 {
				allSingle = false
				break
			}
			flat = append(flat, r[0])
		}
		if allSingle && len(flat)%cols == 0 {
			chunked := make([][]string, 0, len(flat)/cols)
			for i := 0; i < len(flat); i += cols {
				chunked = append(chunked, append([]string{}, flat[i:i+cols]...))
			}
			rows = chunked
		}
	}

	nRows := make([][]string, 0, len(rows))
	for _, r := range rows {
		row := make([]string, cols)
		for i := 0; i < cols; i++ {
			if i < len(r) {
				row[i] = strings.TrimSpace(r[i])
			}
		}
		hasValue := false
		for _, c := range row {
			if c != "" {
				hasValue = true
				break
			}
		}
		if hasValue {
			nRows = append(nRows, row)
		}
	}
	return nHeaders, nRows
}

// safeTranslate calls gtranslate with panic recovery. Returns (result, didFallback).
// On empty input: ("", false). On error/panic/empty result: (trimmed original, true).
// gtranslate scrapes Google Translate's free endpoint and can panic when Google
// returns a rate-limit HTML page instead of JSON; recover() prevents the handler crash.
func safeTranslate(text, fromLang, toLang string) (string, bool) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", false
	}

	var translated string
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("gtranslate panic: %v", r)
			}
		}()
		translated, err = gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{
			From: fromLang,
			To:   toLang,
		})
	}()

	if err != nil || strings.TrimSpace(translated) == "" {
		if err != nil {
			log.Printf("history: translate %s->%s fallback to source: %v", fromLang, toLang, err)
		}
		return trimmed, true
	}
	return strings.TrimSpace(translated), false
}

func translateTextPtrHistory(text *string, fromLang, toLang string) (*string, bool) {
	if text == nil {
		return nil, false
	}
	trimmed := strings.TrimSpace(*text)
	if trimmed == "" {
		return nil, false
	}
	result, fallback := safeTranslate(trimmed, fromLang, toLang)
	return &result, fallback
}

func translateMatrix(headers []string, rows [][]string, fromLang, toLang string) ([]string, [][]string, bool) {
	anyFallback := false
	enHeaders := make([]string, 0, len(headers))
	for _, h := range headers {
		v, fb := safeTranslate(h, fromLang, toLang)
		if fb && strings.TrimSpace(h) != "" {
			anyFallback = true
		}
		enHeaders = append(enHeaders, v)
	}
	enRows := make([][]string, 0, len(rows))
	for _, r := range rows {
		enRow := make([]string, 0, len(r))
		for _, c := range r {
			v, fb := safeTranslate(c, fromLang, toLang)
			if fb && strings.TrimSpace(c) != "" {
				anyFallback = true
			}
			enRow = append(enRow, v)
		}
		enRows = append(enRows, enRow)
	}
	return enHeaders, enRows, anyFallback
}

func toPgText2D(rows [][]string) pgtype.TextArray {
	if len(rows) == 0 {
		return pgtype.TextArray{Status: pgtype.Null}
	}

	cols := 0
	for _, r := range rows {
		if len(r) > cols {
			cols = len(r)
		}
	}
	if cols == 0 {
		return pgtype.TextArray{Status: pgtype.Null}
	}

	elements := make([]pgtype.Text, 0, len(rows)*cols)
	for _, r := range rows {
		for i := 0; i < cols; i++ {
			cell := ""
			if i < len(r) {
				cell = r[i]
			}
			elements = append(elements, pgtype.Text{String: cell, Status: pgtype.Present})
		}
	}

	return pgtype.TextArray{
		Elements: elements,
		Dimensions: []pgtype.ArrayDimension{
			{Length: int32(len(rows)), LowerBound: 1},
			{Length: int32(cols), LowerBound: 1},
		},
		Status: pgtype.Present,
	}
}
