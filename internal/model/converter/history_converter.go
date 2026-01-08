package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"

	"github.com/jackc/pgtype"
)

func HistoryToModel(h entity.History) *model.History {
	return &model.History{
		ID:           h.ID,
		Year:         h.Year,
		Title:        h.Title,
		Description:  h.Description,
		TableHeaders: []string(h.TableHeaders),
		TableRows:    parse2DArray(h.TableRows),
	}
}

func HistoryListToModel(histories []entity.History) []model.History {
	result := make([]model.History, 0, len(histories))

	for _, h := range histories {
		result = append(result, *HistoryToModel(h))
	}

	return result
}

func parse2DArray(arr pgtype.TextArray) [][]string {
	if arr.Status != pgtype.Present {
		return nil
	}

	rows := make([][]string, 0)

	cols := int(arr.Dimensions[1].Length)

	for i := 0; i < len(arr.Elements); i += cols {
		row := make([]string, 0, cols)
		for j := 0; j < cols; j++ {
			row = append(row, arr.Elements[i+j].String)
		}
		rows = append(rows, row)
	}

	return rows
}
