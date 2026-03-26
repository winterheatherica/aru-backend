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

type HeroUsecase interface {
	List(ctx context.Context) ([]model.HeroAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.HeroAdminItem, error)
	Create(ctx context.Context, input model.HeroUpsertInput, image *multipart.FileHeader) (*model.HeroAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.HeroUpsertInput, image *multipart.FileHeader) (*model.HeroAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type heroUsecaseImpl struct {
	heroRepo      repository.HeroRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewHeroUsecase(
	heroRepo repository.HeroRepository,
	minioClient *minio.Client,
	minioBucket string,
	publicBaseURL string,
) HeroUsecase {
	return &heroUsecaseImpl{
		heroRepo:      heroRepo,
		minioClient:   minioClient,
		minioBucket:   minioBucket,
		publicBaseURL: strings.TrimSuffix(publicBaseURL, "/"),
	}
}

func (u *heroUsecaseImpl) List(ctx context.Context) ([]model.HeroAdminItem, error) {
	slides, err := u.heroRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]model.HeroAdminItem, 0, len(slides))
	for _, slide := range slides {
		result = append(result, u.toAdminItem(slide))
	}
	return result, nil
}

func (u *heroUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.HeroAdminItem, error) {
	slide, err := u.heroRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item := u.toAdminItem(*slide)
	return &item, nil
}

func (u *heroUsecaseImpl) Create(ctx context.Context, input model.HeroUpsertInput, image *multipart.FileHeader) (*model.HeroAdminItem, error) {
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

	translations := []entity.HeroSlideTranslation{
		{
			Language:    strings.ToUpper(input.Language),
			Alt:         input.Alt,
			Title:       input.Title,
			Description: input.Description,
		},
	}

	if strings.EqualFold(input.Language, "ID") {
		translations = append(translations, entity.HeroSlideTranslation{
			Language:    "EN",
			Alt:         translateTextPtr(input.Alt, "id", "en"),
			Title:       translateTextPtr(input.Title, "id", "en"),
			Description: translateTextPtr(input.Description, "id", "en"),
		})
	}

	slide := &entity.HeroSlide{
		ImagePath:    objectName,
		OrderIndex:   input.OrderIndex,
		IsActive:     input.IsActive,
		Banner:       input.Banner,
		Translations: translations,
	}

	if err := u.heroRepo.Create(ctx, slide); err != nil {
		_ = u.removeObject(ctx, objectName)
		return nil, err
	}

	created, err := u.heroRepo.FindByID(ctx, slide.ID)
	if err != nil {
		item := u.toAdminItem(*slide)
		return &item, nil
	}
	item := u.toAdminItem(*created)
	return &item, nil
}

func (u *heroUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.HeroUpsertInput, image *multipart.FileHeader) (*model.HeroAdminItem, error) {
	slide, err := u.heroRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input.Language) == "" {
		input.Language = "ID"
	}

	oldObject := slide.ImagePath
	if image != nil {
		newObject, upErr := u.uploadImage(ctx, image)
		if upErr != nil {
			return nil, upErr
		}
		slide.ImagePath = newObject
	}

	slide.OrderIndex = input.OrderIndex
	slide.IsActive = input.IsActive
	slide.Banner = input.Banner

	if err := u.heroRepo.Update(ctx, slide); err != nil {
		return nil, err
	}

	tr := &entity.HeroSlideTranslation{
		HeroSlideID: slide.ID,
		Language:    strings.ToUpper(input.Language),
		Alt:         input.Alt,
		Title:       input.Title,
		Description: input.Description,
	}
	if err := u.heroRepo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	if image != nil && oldObject != "" && oldObject != slide.ImagePath {
		_ = u.removeObject(ctx, oldObject)
	}

	updated, err := u.heroRepo.FindByID(ctx, slide.ID)
	if err != nil {
		item := u.toAdminItem(*slide)
		return &item, nil
	}
	item := u.toAdminItem(*updated)
	return &item, nil
}

func (u *heroUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	slide, err := u.heroRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.heroRepo.DeleteByID(ctx, id); err != nil {
		return err
	}

	if slide.ImagePath != "" {
		_ = u.removeObject(ctx, slide.ImagePath)
	}

	return nil
}

func (u *heroUsecaseImpl) toAdminItem(slide entity.HeroSlide) model.HeroAdminItem {
	translations := make([]model.HeroAdminTranslation, 0, len(slide.Translations))
	for _, tr := range slide.Translations {
		translations = append(translations, model.HeroAdminTranslation{
			Language:    tr.Language,
			Alt:         tr.Alt,
			Title:       tr.Title,
			Description: tr.Description,
		})
	}

	return model.HeroAdminItem{
		ID:           slide.ID,
		ImagePath:    slide.ImagePath,
		MainImageURL: converter.BuildAssetURL(u.publicBaseURL, slide.ImagePath),
		OrderIndex:   slide.OrderIndex,
		IsActive:     slide.IsActive,
		Banner:       slide.Banner,
		Translations: translations,
	}
}

func translateTextPtr(text *string, fromLang, toLang string) *string {
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

func (u *heroUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/hero_slides/%s%s", uuid.NewString(), ext)

	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{
		ContentType: image.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (u *heroUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}

func normalizeObjectName(path, publicBaseURL, bucket string) string {
	p := strings.TrimSpace(path)
	if p == "" {
		return ""
	}

	p = strings.TrimPrefix(p, publicBaseURL)
	p = strings.TrimPrefix(p, "/")
	p = strings.TrimPrefix(p, bucket+"/")
	p = strings.TrimPrefix(p, "api/")
	p = strings.TrimPrefix(p, "assets/")
	return p
}
