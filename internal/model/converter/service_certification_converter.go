package converter

import (
	"aru-backend/internal/entity"
	"aru-backend/internal/model"
)

func ServiceCertificationToModel(
	cert entity.ServiceCertification,
	lang string,
) *model.ServiceCertification {

	tr := findServiceCertificationTranslation(cert, lang)
	if tr == nil {
		return nil
	}

	return &model.ServiceCertification{
		ID: cert.ID,

		Service: cert.Service,

		OrderIndex: cert.OrderIndex,
		IsActive:   cert.IsActive,

		Title:   tr.Title,
		Alt:     tr.Alt,
		Caption: tr.Caption,
	}
}

func findServiceCertificationTranslation(
	cert entity.ServiceCertification,
	lang string,
) *entity.ServiceCertificationTranslation {

	for i := range cert.Translations {
		if cert.Translations[i].Language == lang {
			return &cert.Translations[i]
		}
	}
	return nil
}

func ServiceCertificationListToModel(
	certs []entity.ServiceCertification,
	lang string,
) []model.ServiceCertification {

	result := make([]model.ServiceCertification, 0, len(certs))

	for _, c := range certs {
		m := ServiceCertificationToModel(c, lang)
		if m != nil {
			result = append(result, *m)
		}
	}

	return result
}
