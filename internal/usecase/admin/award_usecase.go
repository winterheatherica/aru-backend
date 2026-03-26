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

type AwardUsecase interface {
	List(ctx context.Context) ([]model.AwardAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.AwardAdminItem, error)
	Create(ctx context.Context, input model.AwardUpsertInput, image *multipart.FileHeader) (*model.AwardAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.AwardUpsertInput, image *multipart.FileHeader) (*model.AwardAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type awardUsecaseImpl struct {
	awardRepo     repository.AwardRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewAwardUsecase(
	awardRepo repository.AwardRepository,
	minioClient *minio.Client,
	minioBucket string,
	publicBaseURL string,
) AwardUsecase {
	return &awardUsecaseImpl{
		awardRepo:     awardRepo,
		minioClient:   minioClient,
		minioBucket:   minioBucket,
		publicBaseURL: strings.TrimSuffix(publicBaseURL, "/"),
	}
}

func (u *awardUsecaseImpl) List(ctx context.Context) ([]model.AwardAdminItem, error) {
	items, err := u.awardRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.AwardAdminItem, 0, len(items))
	for _, item := range items {
		result = append(result, u.toAdminItem(item))
	}
	return result, nil
}

func (u *awardUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.AwardAdminItem, error) {
	item, err := u.awardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*item)
	return &res, nil
}

func (u *awardUsecaseImpl) Create(ctx context.Context, input model.AwardUpsertInput, image *multipart.FileHeader) (*model.AwardAdminItem, error) {
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

	translations := []entity.AwardTranslation{
		{
			Language:    strings.ToUpper(input.Language),
			Alt:         input.Alt,
			Title:       input.Title,
			Label:       input.Label,
			Description: input.Description,
		},
	}

	if strings.EqualFold(input.Language, "ID") {
		translations = append(translations, entity.AwardTranslation{
			Language:    "EN",
			Alt:         translateTextPtrAward(input.Alt, "id", "en"),
			Title:       translateTextPtrAward(input.Title, "id", "en"),
			Label:       translateTextPtrAward(input.Label, "id", "en"),
			Description: translateTextPtrAward(input.Description, "id", "en"),
		})
	}

	item := &entity.Award{
		ImagePath:    objectName,
		Year:         input.Year,
		OrderIndex:   input.OrderIndex,
		IsActive:     input.IsActive,
		Translations: translations,
	}

	if err := u.awardRepo.Create(ctx, item); err != nil {
		_ = u.removeObject(ctx, objectName)
		return nil, err
	}

	created, err := u.awardRepo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item)
		return &res, nil
	}
	res := u.toAdminItem(*created)
	return &res, nil
}

func (u *awardUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.AwardUpsertInput, image *multipart.FileHeader) (*model.AwardAdminItem, error) {
	item, err := u.awardRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	oldObject := item.ImagePath
	if image != nil {
		newObject, upErr := u.uploadImage(ctx, image)
		if upErr != nil {
			return nil, upErr
		}
		item.ImagePath = newObject
	}

	item.Year = input.Year
	item.OrderIndex = input.OrderIndex
	item.IsActive = input.IsActive

	if err := u.awardRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.AwardTranslation{
		AwardID:     item.ID,
		Language:    strings.ToUpper(input.Language),
		Alt:         input.Alt,
		Title:       input.Title,
		Label:       input.Label,
		Description: input.Description,
	}
	if err := u.awardRepo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	if image != nil && oldObject != "" && oldObject != item.ImagePath {
		_ = u.removeObject(ctx, oldObject)
	}

	updated, err := u.awardRepo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item)
		return &res, nil
	}
	res := u.toAdminItem(*updated)
	return &res, nil
}

func (u *awardUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	item, err := u.awardRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.awardRepo.DeleteByID(ctx, id); err != nil {
		return err
	}

	if item.ImagePath != "" {
		_ = u.removeObject(ctx, item.ImagePath)
	}

	return nil
}

func (u *awardUsecaseImpl) toAdminItem(item entity.Award) model.AwardAdminItem {
	translations := make([]model.AwardAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		translations = append(translations, model.AwardAdminTranslation{
			Language:    tr.Language,
			Alt:         tr.Alt,
			Title:       tr.Title,
			Label:       tr.Label,
			Description: tr.Description,
		})
	}

	return model.AwardAdminItem{
		ID:           item.ID,
		ImagePath:    item.ImagePath,
		ImageURL:     converter.BuildAssetURL(u.publicBaseURL, item.ImagePath),
		Year:         item.Year,
		OrderIndex:   item.OrderIndex,
		IsActive:     item.IsActive,
		Translations: translations,
	}
}

func translateTextPtrAward(text *string, fromLang, toLang string) *string {
	if text == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*text)
	if trimmed == "" {
		return nil
	}

	translated, err := gtranslate.TranslateWithParams(trimmed, gtranslate.TranslationParams{
		From: fromLang,
		To:   toLang,
	})
	if err != nil || strings.TrimSpace(translated) == "" {
		fallback := trimmed
		return &fallback
	}

	result := strings.TrimSpace(translated)
	return &result
}

func (u *awardUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/awards/%s%s", uuid.NewString(), ext)

	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{
		ContentType: image.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *awardUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
