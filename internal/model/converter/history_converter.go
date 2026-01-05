package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func HistoryToModel(h entity.History) *model.History {
	return &model.History{
		ID:           h.ID,
		Year:         h.Year,
		Title:        h.Title,
		Description:  h.Description,
		TableHeaders: h.TableHeaders,
		TableRows:    h.TableRows,
	}
}

func HistoryListToModel(histories []entity.History) []model.History {
	result := make([]model.History, 0, len(histories))

	for _, h := range histories {
		result = append(result, *HistoryToModel(h))
	}

	return result
}
