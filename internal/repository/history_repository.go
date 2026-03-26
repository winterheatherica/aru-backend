package repository

import (
	"context"
	"fmt"
	"strings"

	"aru-backend/internal/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type HistoryRepository interface {
	FindActiveByLanguage(ctx context.Context, lang string) ([]entity.History, error)
	FindAll(ctx context.Context) ([]entity.History, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.History, error)
	FindByYearAndLanguage(ctx context.Context, year int, lang string) (*entity.History, error)
	Create(ctx context.Context, item *entity.History) error
	Update(ctx context.Context, item *entity.History) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type historyRepositoryImpl struct{ db *gorm.DB }

func NewHistoryRepository(db *gorm.DB) HistoryRepository { return &historyRepositoryImpl{db: db} }

func (r *historyRepositoryImpl) FindActiveByLanguage(ctx context.Context, lang string) ([]entity.History, error) {
	var histories []entity.History
	err := r.db.WithContext(ctx).
		Where("language = ?", lang).
		Where("is_active = ?", true).
		Order("year DESC NULLS LAST").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *historyRepositoryImpl) FindAll(ctx context.Context) ([]entity.History, error) {
	var histories []entity.History
	if err := r.db.WithContext(ctx).
		Order("year DESC NULLS LAST, language ASC, created_at DESC").
		Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *historyRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*entity.History, error) {
	var item entity.History
	if err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *historyRepositoryImpl) FindByYearAndLanguage(ctx context.Context, year int, lang string) (*entity.History, error) {
	var item entity.History
	if err := r.db.WithContext(ctx).
		Where("year = ? AND language = ?", year, lang).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *historyRepositoryImpl) Create(ctx context.Context, item *entity.History) error {
	rowsLiteral := textArray2DLiteral(item.TableRows)

	var idStr string
	err := r.db.WithContext(ctx).Raw(`
		INSERT INTO histories (language, year, title, description, table_headers, table_rows, is_active, uploaded_by)
		VALUES (?, ?, ?, ?, ?, ?::text[][], ?, ?)
		RETURNING id::text
	`,
		item.Language,
		item.Year,
		item.Title,
		item.Description,
		pq.Array([]string(item.TableHeaders)),
		rowsLiteral,
		item.IsActive,
		item.UploadedBy,
	).Scan(&idStr).Error
	if err != nil {
		return err
	}
	parsed, perr := uuid.Parse(idStr)
	if perr != nil {
		return perr
	}
	item.ID = parsed
	return nil
}

func (r *historyRepositoryImpl) Update(ctx context.Context, item *entity.History) error {
	rowsLiteral := textArray2DLiteral(item.TableRows)

	return r.db.WithContext(ctx).Exec(`
		UPDATE histories
		SET language = ?,
			year = ?,
			title = ?,
			description = ?,
			table_headers = ?,
			table_rows = ?::text[][],
			is_active = ?
		WHERE id = ?
	`,
		item.Language,
		item.Year,
		item.Title,
		item.Description,
		pq.Array([]string(item.TableHeaders)),
		rowsLiteral,
		item.IsActive,
		item.ID,
	).Error
}

func (r *historyRepositoryImpl) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entity.History{}, "id = ?", id).Error
}

func textArray2DLiteral(arr pgtype.TextArray) any {
	if arr.Status != pgtype.Present || len(arr.Dimensions) < 2 {
		return nil
	}

	rows := int(arr.Dimensions[0].Length)
	cols := int(arr.Dimensions[1].Length)
	if rows <= 0 || cols <= 0 {
		return nil
	}

	var b strings.Builder
	b.WriteString("{")
	idx := 0
	for r := 0; r < rows; r++ {
		if r > 0 {
			b.WriteString(",")
		}
		b.WriteString("{")
		for c := 0; c < cols; c++ {
			if c > 0 {
				b.WriteString(",")
			}
			val := ""
			if idx < len(arr.Elements) {
				val = arr.Elements[idx].String
			}
			idx++
			val = strings.ReplaceAll(val, `\`, `\\`)
			val = strings.ReplaceAll(val, `"`, `\"`)
			b.WriteString(fmt.Sprintf(`"%s"`, val))
		}
		b.WriteString("}")
	}
	b.WriteString("}")
	return b.String()
}
