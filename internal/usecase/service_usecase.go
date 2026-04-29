package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type ServiceUsecase interface {
	GetServicePage(ctx context.Context, service string, lang string) (*ServiceResponse, error)
	GetServicePageIsolated(ctx context.Context, service string, lang string) (*ServiceResponse, map[string]string)
}

type serviceUsecaseImpl struct {
	galleryRepo       repository.ServiceGalleryRepository
	pricingRepo       repository.ServicePricingRepository
	matrixRepo        repository.ServiceMatrixRepository
	certificationRepo repository.ServiceCertificationRepository

	baseURL string
}

type ServiceResponse struct {
	Gallery        []model.ServiceGallery       `json:"gallery"`
	Pricing        []model.ServicePricingTier   `json:"pricing"`
	Matrix         *model.ServiceMatrix         `json:"matrix"`
	Certifications []model.ServiceCertification `json:"certifications"`
}

func NewServiceUsecase(
	galleryRepo repository.ServiceGalleryRepository,
	pricingRepo repository.ServicePricingRepository,
	matrixRepo repository.ServiceMatrixRepository,
	certificationRepo repository.ServiceCertificationRepository,
	baseURL string,
) ServiceUsecase {
	return &serviceUsecaseImpl{
		galleryRepo:       galleryRepo,
		pricingRepo:       pricingRepo,
		matrixRepo:        matrixRepo,
		certificationRepo: certificationRepo,
		baseURL:           baseURL,
	}
}

func (u *serviceUsecaseImpl) GetServicePage(
	ctx context.Context,
	service string,
	lang string,
) (*ServiceResponse, error) {
	galleryEntities, err := u.galleryRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		return nil, err
	}
	gallery := converter.ServiceGalleryListToModel(galleryEntities, lang, u.baseURL)

	pricingEntities, err := u.pricingRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		return nil, err
	}
	pricing := converter.ServicePricingTierListToModel(pricingEntities, lang)

	matrixEntity, err := u.matrixRepo.FindActiveByService(ctx, service, nil, lang)
	if err != nil {
		return nil, err
	}

	var matrix *model.ServiceMatrix
	if matrixEntity != nil {
		matrix = converter.ServiceMatrixToModel(*matrixEntity, lang)
	}

	certEntities, err := u.certificationRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		return nil, err
	}
	certifications := converter.ServiceCertificationListToModel(certEntities, lang)

	return &ServiceResponse{Gallery: gallery, Pricing: pricing, Matrix: matrix, Certifications: certifications}, nil
}

func (u *serviceUsecaseImpl) GetServicePageIsolated(
	ctx context.Context,
	service string,
	lang string,
) (*ServiceResponse, map[string]string) {
	resp := &ServiceResponse{}
	sectionErrors := map[string]string{}

	galleryEntities, err := u.galleryRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		sectionErrors["gallery"] = err.Error()
		resp.Gallery = []model.ServiceGallery{}
	} else {
		resp.Gallery = converter.ServiceGalleryListToModel(galleryEntities, lang, u.baseURL)
	}

	pricingEntities, err := u.pricingRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		sectionErrors["pricing"] = err.Error()
		resp.Pricing = []model.ServicePricingTier{}
	} else {
		resp.Pricing = converter.ServicePricingTierListToModel(pricingEntities, lang)
	}

	matrixEntity, err := u.matrixRepo.FindActiveByService(ctx, service, nil, lang)
	if err != nil {
		sectionErrors["matrix"] = err.Error()
		resp.Matrix = nil
	} else if matrixEntity != nil {
		resp.Matrix = converter.ServiceMatrixToModel(*matrixEntity, lang)
	}

	certEntities, err := u.certificationRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		sectionErrors["certifications"] = err.Error()
		resp.Certifications = []model.ServiceCertification{}
	} else {
		resp.Certifications = converter.ServiceCertificationListToModel(certEntities, lang)
	}

	return resp, sectionErrors
}
