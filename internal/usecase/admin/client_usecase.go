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

type ClientUsecase interface {
	List(ctx context.Context) ([]model.ClientAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.ClientAdminItem, error)
	Create(ctx context.Context, input model.ClientUpsertInput, image *multipart.FileHeader) (*model.ClientAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.ClientUpsertInput, image *multipart.FileHeader) (*model.ClientAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type clientUsecaseImpl struct {
	repo          repository.ClientRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewClientUsecase(repo repository.ClientRepository, minioClient *minio.Client, minioBucket string, publicBaseURL string) ClientUsecase {
	return &clientUsecaseImpl{repo: repo, minioClient: minioClient, minioBucket: minioBucket, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

func (u *clientUsecaseImpl) List(ctx context.Context) ([]model.ClientAdminItem, error) {
	items, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.ClientAdminItem, 0, len(items))
	for _, it := range items {
		out = append(out, u.toAdminItem(it))
	}
	return out, nil
}

func (u *clientUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.ClientAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*it)
	return &res, nil
}

func (u *clientUsecaseImpl) Create(ctx context.Context, input model.ClientUpsertInput, image *multipart.FileHeader) (*model.ClientAdminItem, error) {
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

	translations := []entity.ClientTranslation{{
		Language:    strings.ToUpper(input.Language),
		Alt:         input.Alt,
		Title:       input.Title,
		Description: input.Description,
	}}
	if strings.EqualFold(input.Language, "ID") {
		translations = append(translations, entity.ClientTranslation{
			Language:    "EN",
			Alt:         translateTextPtrClient(input.Alt, "id", "en"),
			Title:       translateTextPtrClient(input.Title, "id", "en"),
			Description: translateTextPtrClient(input.Description, "id", "en"),
		})
	}

	item := &entity.Client{ImagePath: objectName, OrderIndex: input.OrderIndex, IsActiveClientScroller: input.IsActiveClientScroller, Translations: translations}
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

func (u *clientUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.ClientUpsertInput, image *multipart.FileHeader) (*model.ClientAdminItem, error) {
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
	item.IsActiveClientScroller = input.IsActiveClientScroller

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}

	tr := &entity.ClientTranslation{ClientID: item.ID, Language: strings.ToUpper(input.Language), Alt: input.Alt, Title: input.Title, Description: input.Description}
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

func (u *clientUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
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

func (u *clientUsecaseImpl) toAdminItem(item entity.Client) model.ClientAdminItem {
	trs := make([]model.ClientAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		trs = append(trs, model.ClientAdminTranslation{Language: tr.Language, Alt: tr.Alt, Title: tr.Title, Description: tr.Description})
	}
	return model.ClientAdminItem{ID: item.ID, ImagePath: item.ImagePath, ImageURL: converter.BuildAssetURL(u.publicBaseURL, item.ImagePath), OrderIndex: item.OrderIndex, IsActiveClientScroller: item.IsActiveClientScroller, Translations: trs}
}

func translateTextPtrClient(text *string, fromLang, toLang string) *string {
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

func (u *clientUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()
	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/clients/%s%s", uuid.NewString(), ext)
	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{ContentType: image.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (u *clientUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
