package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func ClientToModel(
	client entity.Client,
	lang string,
	baseURL string,
) *model.Client {

	trans := findClientTranslation(client, lang)

	alt := ""
	var title *string
	var desc *string

	if trans != nil {
		if trans.Alt != nil {
			alt = *trans.Alt
		}
		title = trans.Title
		desc = trans.Description
	}

	return &model.Client{
		ID:    client.ID,
		Src:   BuildAssetURL(baseURL, client.ImagePath),
		Alt:   alt,
		Title: title,
		Desc:  desc,
		Order: client.OrderIndex,
	}
}

func ClientListToModel(
	clients []entity.Client,
	lang string,
	baseURL string,
) []model.Client {

	result := make([]model.Client, 0, len(clients))

	for _, c := range clients {
		result = append(result, *ClientToModel(c, lang, baseURL))
	}

	return result
}

func findClientTranslation(
	client entity.Client,
	lang string,
) *entity.ClientTranslation {

	for i := range client.Translations {
		if client.Translations[i].Language == lang {
			return &client.Translations[i]
		}
	}
	return nil
}
