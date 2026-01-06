package usecase

import (
	"context"

	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"
)

type ServiceUsecase interface {
	GetServicePage(ctx context.Context, service string, lang string) (*ServiceResponse, error)
}

type serviceUsecaseImpl struct {
	galleryRepo repository.ServiceGalleryRepository
	pricingRepo repository.ServicePricingRepository
	matrixRepo  repository.ServiceMatrixRepository

	baseURL string
}

type ServiceResponse struct {
	Gallery []model.ServiceGallery     `json:"gallery"`
	Pricing []model.ServicePricingTier `json:"pricing"`
	Matrix  *model.ServiceMatrix       `json:"matrix"`
}

func NewServiceUsecase(
	galleryRepo repository.ServiceGalleryRepository,
	pricingRepo repository.ServicePricingRepository,
	matrixRepo repository.ServiceMatrixRepository,
	baseURL string,
) ServiceUsecase {
	return &serviceUsecaseImpl{
		galleryRepo: galleryRepo,
		pricingRepo: pricingRepo,
		matrixRepo:  matrixRepo,
		baseURL:     baseURL,
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
	gallery := converter.ServiceGalleryListToModel(
		galleryEntities,
		lang,
		u.baseURL,
	)

	pricingEntities, err := u.pricingRepo.FindActiveByService(ctx, service, lang)
	if err != nil {
		return nil, err
	}
	pricing := converter.ServicePricingTierListToModel(
		pricingEntities,
		lang,
	)

	matrixEntity, err := u.matrixRepo.FindActiveByService(ctx, service, nil, lang)
	if err != nil {
		return nil, err
	}
	matrix := converter.ServiceMatrixToModel(
		*matrixEntity,
		lang,
	)

	return &ServiceResponse{
		Gallery: gallery,
		Pricing: pricing,
		Matrix:  matrix,
	}, nil
}
