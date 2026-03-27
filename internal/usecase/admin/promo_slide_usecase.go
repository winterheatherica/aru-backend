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

type PromoSlideUsecase interface {
	List(ctx context.Context) ([]model.PromoSlideAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.PromoSlideAdminItem, error)
	Create(ctx context.Context, input model.PromoSlideUpsertInput, image *multipart.FileHeader) (*model.PromoSlideAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.PromoSlideUpsertInput, image *multipart.FileHeader) (*model.PromoSlideAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type promoSlideUsecaseImpl struct {
	repo          repository.PromoRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewPromoSlideUsecase(repo repository.PromoRepository, minioClient *minio.Client, minioBucket string, publicBaseURL string) PromoSlideUsecase {
	return &promoSlideUsecaseImpl{repo: repo, minioClient: minioClient, minioBucket: minioBucket, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

func (u *promoSlideUsecaseImpl) List(ctx context.Context) ([]model.PromoSlideAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.PromoSlideAdminItem, 0, len(items))
	for _, it := range items {
		out = append(out, u.toAdminItem(it))
	}
	return out, nil
}

func (u *promoSlideUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.PromoSlideAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*it)
	return &res, nil
}

func (u *promoSlideUsecaseImpl) Create(ctx context.Context, input model.PromoSlideUpsertInput, image *multipart.FileHeader) (*model.PromoSlideAdminItem, error) {
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

	translations := []entity.PromoSlideTranslation{{
		Language: strings.ToUpper(input.Language),
		Alt:      input.Alt,
		Title:    input.Title,
	}}
	if strings.EqualFold(input.Language, "ID") {
		translations = append(translations, entity.PromoSlideTranslation{
			Language: "EN",
			Alt:      translateTextPtrPromo(input.Alt, "id", "en"),
			Title:    translateTextPtrPromo(input.Title, "id", "en"),
		})
	}

	item := &entity.PromoSlide{ImagePath: objectName, OrderIndex: input.OrderIndex, IsActive: input.IsActive, Translations: translations}
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

func (u *promoSlideUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.PromoSlideUpsertInput, image *multipart.FileHeader) (*model.PromoSlideAdminItem, error) {
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
	item.IsActive = input.IsActive

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.PromoSlideTranslation{PromoSlideID: item.ID, Language: strings.ToUpper(input.Language), Alt: input.Alt, Title: input.Title}
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

func (u *promoSlideUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
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

func (u *promoSlideUsecaseImpl) toAdminItem(item entity.PromoSlide) model.PromoSlideAdminItem {
	trs := make([]model.PromoSlideAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		trs = append(trs, model.PromoSlideAdminTranslation{Language: tr.Language, Alt: tr.Alt, Title: tr.Title})
	}
	return model.PromoSlideAdminItem{ID: item.ID, ImagePath: item.ImagePath, ImageURL: converter.BuildAssetURL(u.publicBaseURL, item.ImagePath), OrderIndex: item.OrderIndex, IsActive: item.IsActive, Translations: trs}
}

func translateTextPtrPromo(text *string, fromLang, toLang string) *string {
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

func (u *promoSlideUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/promo_slides/%s%s", uuid.NewString(), ext)
	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{ContentType: image.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (u *promoSlideUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
