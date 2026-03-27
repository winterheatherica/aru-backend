package admin

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type PartnerUsecase interface {
	List(ctx context.Context) ([]model.PartnerAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.PartnerAdminItem, error)
	Create(ctx context.Context, input model.PartnerUpsertInput, image *multipart.FileHeader) (*model.PartnerAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.PartnerUpsertInput, image *multipart.FileHeader) (*model.PartnerAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type partnerUsecaseImpl struct {
	repo          repository.PartnerRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewPartnerUsecase(repo repository.PartnerRepository, minioClient *minio.Client, minioBucket string, publicBaseURL string) PartnerUsecase {
	return &partnerUsecaseImpl{repo: repo, minioClient: minioClient, minioBucket: minioBucket, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

func (u *partnerUsecaseImpl) List(ctx context.Context) ([]model.PartnerAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.PartnerAdminItem, 0, len(items))
	for _, it := range items {
		out = append(out, u.toAdminItem(it))
	}
	return out, nil
}

func (u *partnerUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.PartnerAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*it)
	return &res, nil
}

func (u *partnerUsecaseImpl) Create(ctx context.Context, input model.PartnerUpsertInput, image *multipart.FileHeader) (*model.PartnerAdminItem, error) {
	if image == nil {
		return nil, fmt.Errorf("image is required")
	}
	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	objectName, err := u.uploadImage(ctx, image)
	if err != nil {
		return nil, err
	}

	translations := []entity.PartnerTranslation{{
		Language:    strings.ToUpper(input.Language),
		Alt:         input.Alt,
		Title:       input.Title,
		Description: input.Description,
	}}
	if strings.EqualFold(input.Language, "ID") {
		translations = append(translations, entity.PartnerTranslation{
			Language:    "EN",
			Alt:         translateTextPtrPartner(input.Alt, "id", "en"),
			Title:       translateTextPtrPartner(input.Title, "id", "en"),
			Description: translateTextPtrPartner(input.Description, "id", "en"),
		})
	}

	item := &entity.Partner{
		ImagePath:               objectName,
		OrderIndex:              input.OrderIndex,
		IsActivePartnerGrid:     input.IsActivePartnerGrid,
		IsActivePartnerScroller: input.IsActivePartnerScroller,
		Translations:            translations,
	}

	if err := u.repo.Create(ctx, item); err != nil {
		_ = u.removeObject(ctx, objectName)
		return nil, err
	}

	created, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item)
		return &res, nil
	}
	res := u.toAdminItem(*created)
	return &res, nil
}

func (u *partnerUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.PartnerUpsertInput, image *multipart.FileHeader) (*model.PartnerAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	oldObject := item.ImagePath
	if image != nil {
		newObj, upErr := u.uploadImage(ctx, image)
		if upErr != nil {
			return nil, upErr
		}
		item.ImagePath = newObj
	}
	item.OrderIndex = input.OrderIndex
	item.IsActivePartnerGrid = input.IsActivePartnerGrid
	item.IsActivePartnerScroller = input.IsActivePartnerScroller

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.PartnerTranslation{PartnerID: item.ID, Language: strings.ToUpper(input.Language), Alt: input.Alt, Title: input.Title, Description: input.Description}
	if err := u.repo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	if image != nil && oldObject != "" && oldObject != item.ImagePath {
		_ = u.removeObject(ctx, oldObject)
	}

	updated, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item)
		return &res, nil
	}
	res := u.toAdminItem(*updated)
	return &res, nil
}

func (u *partnerUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := u.repo.DeleteByID(ctx, id); err != nil {
		return err
	}
	if item.ImagePath != "" {
		_ = u.removeObject(ctx, item.ImagePath)
	}
	return nil
}

func (u *partnerUsecaseImpl) toAdminItem(item entity.Partner) model.PartnerAdminItem {
	trs := make([]model.PartnerAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		trs = append(trs, model.PartnerAdminTranslation{Language: tr.Language, Alt: tr.Alt, Title: tr.Title, Description: tr.Description})
	}
	return model.PartnerAdminItem{
		ID:                      item.ID,
		ImagePath:               item.ImagePath,
		ImageURL:                converter.BuildAssetURL(u.publicBaseURL, item.ImagePath),
		OrderIndex:              item.OrderIndex,
		IsActivePartnerGrid:     item.IsActivePartnerGrid,
		IsActivePartnerScroller: item.IsActivePartnerScroller,
		Translations:            trs,
	}
}

func translateTextPtrPartner(text *string, fromLang, toLang string) *string {
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

func (u *partnerUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/partners/%s%s", uuid.NewString(), ext)
	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{ContentType: image.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (u *partnerUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
