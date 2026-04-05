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

type ServiceGalleryUsecase interface {
	ListByService(ctx context.Context, service string, lang string) ([]model.ServiceGalleryAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceGalleryAdminItem, error)
	Create(ctx context.Context, input model.ServiceGalleryUpsertInput, image *multipart.FileHeader) (*model.ServiceGalleryAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.ServiceGalleryUpsertInput, image *multipart.FileHeader) (*model.ServiceGalleryAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type serviceGalleryUsecaseImpl struct {
	repo          repository.ServiceGalleryRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewServiceGalleryUsecase(repo repository.ServiceGalleryRepository, minioClient *minio.Client, minioBucket string, publicBaseURL string) ServiceGalleryUsecase {
	return &serviceGalleryUsecaseImpl{repo: repo, minioClient: minioClient, minioBucket: minioBucket, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

func (u *serviceGalleryUsecaseImpl) ListByService(ctx context.Context, service string, lang string) ([]model.ServiceGalleryAdminItem, error) {
	items, err := u.repo.FindByService(ctx, strings.ToUpper(strings.TrimSpace(service)))
	if err != nil {
		return nil, err
	}
	res := make([]model.ServiceGalleryAdminItem, 0, len(items))
	for _, it := range items {
		res = append(res, u.toAdminItem(it, lang))
	}
	return res, nil
}

func (u *serviceGalleryUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID, lang string) (*model.ServiceGalleryAdminItem, error) {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	res := u.toAdminItem(*it, lang)
	return &res, nil
}

func (u *serviceGalleryUsecaseImpl) Create(ctx context.Context, input model.ServiceGalleryUpsertInput, image *multipart.FileHeader) (*model.ServiceGalleryAdminItem, error) {
	if image == nil {
		return nil, fmt.Errorf("image is required")
	}
	lang := normServiceGalleryLang(input.Language)
	service := strings.ToUpper(strings.TrimSpace(input.Service))
	mediaType := strings.ToUpper(strings.TrimSpace(input.MediaType))
	if mediaType == "" {
		mediaType = "IMAGE"
	}

	objectName, err := u.uploadImage(ctx, image)
	if err != nil {
		return nil, err
	}

	item := &entity.ServiceGallery{Service: service, MediaType: mediaType, Src: objectName, IsActive: input.IsActive}
	if err := u.repo.Create(ctx, item); err != nil {
		_ = u.removeObject(ctx, objectName)
		return nil, err
	}

	if err := u.repo.UpsertTranslation(ctx, &entity.ServiceGalleryTranslation{GalleryID: item.ID, Language: lang, Title: trimPtr(input.Title), Alt: trimPtr(input.Alt), Caption: trimPtr(input.Caption)}); err != nil {
		return nil, err
	}
	if lang == "ID" {
		_ = u.repo.UpsertTranslation(ctx, &entity.ServiceGalleryTranslation{GalleryID: item.ID, Language: "EN", Title: translateTextPtrServiceGallery(input.Title, "id", "en"), Alt: translateTextPtrServiceGallery(input.Alt, "id", "en"), Caption: translateTextPtrServiceGallery(input.Caption, "id", "en")})
	}

	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item, lang)
		return &res, nil
	}
	res := u.toAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceGalleryUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.ServiceGalleryUpsertInput, image *multipart.FileHeader) (*model.ServiceGalleryAdminItem, error) {
	item, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	lang := normServiceGalleryLang(input.Language)
	if s := strings.ToUpper(strings.TrimSpace(input.Service)); s != "" {
		item.Service = s
	}
	if mt := strings.ToUpper(strings.TrimSpace(input.MediaType)); mt != "" {
		item.MediaType = mt
	}
	item.IsActive = input.IsActive

	oldObject := item.Src
	if image != nil {
		newObj, upErr := u.uploadImage(ctx, image)
		if upErr != nil {
			return nil, upErr
		}
		item.Src = newObj
	}

	if err := u.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	if err := u.repo.UpsertTranslation(ctx, &entity.ServiceGalleryTranslation{GalleryID: item.ID, Language: lang, Title: trimPtr(input.Title), Alt: trimPtr(input.Alt), Caption: trimPtr(input.Caption)}); err != nil {
		return nil, err
	}
	if lang == "ID" {
		_ = u.repo.UpsertTranslation(ctx, &entity.ServiceGalleryTranslation{GalleryID: item.ID, Language: "EN", Title: translateTextPtrServiceGallery(input.Title, "id", "en"), Alt: translateTextPtrServiceGallery(input.Alt, "id", "en"), Caption: translateTextPtrServiceGallery(input.Caption, "id", "en")})
	}
	if image != nil && oldObject != "" && oldObject != item.Src {
		_ = u.removeObject(ctx, oldObject)
	}

	fresh, err := u.repo.FindByID(ctx, item.ID)
	if err != nil {
		res := u.toAdminItem(*item, lang)
		return &res, nil
	}
	res := u.toAdminItem(*fresh, lang)
	return &res, nil
}

func (u *serviceGalleryUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error {
	it, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := u.repo.DeleteByID(ctx, id); err != nil {
		return err
	}
	if it.Src != "" {
		_ = u.removeObject(ctx, it.Src)
	}
	return nil
}

func (u *serviceGalleryUsecaseImpl) toAdminItem(item entity.ServiceGallery, lang string) model.ServiceGalleryAdminItem {
	trs := make([]model.ServiceGalleryAdminTranslation, 0, len(item.Translations))
	for _, tr := range item.Translations {
		trs = append(trs, model.ServiceGalleryAdminTranslation{Language: tr.Language, Title: tr.Title, Alt: tr.Alt, Caption: tr.Caption})
	}
	var title, alt, caption *string
	for _, tr := range item.Translations {
		if strings.EqualFold(tr.Language, lang) {
			title, alt, caption = tr.Title, tr.Alt, tr.Caption
			break
		}
	}
	if title == nil && len(item.Translations) > 0 {
		title, alt, caption = item.Translations[0].Title, item.Translations[0].Alt, item.Translations[0].Caption
	}
	img := converter.BuildAssetURL(u.publicBaseURL, item.Src)
	var thumb *string
	if item.Thumbnail != nil {
		t := converter.BuildAssetURL(u.publicBaseURL, *item.Thumbnail)
		thumb = &t
	}
	return model.ServiceGalleryAdminItem{ID: item.ID, Service: item.Service, MediaType: item.MediaType, Src: item.Src, ImageURL: img, Thumbnail: thumb, IsActive: item.IsActive, Title: title, Alt: alt, Caption: caption, Translations: trs}
}

func normServiceGalleryLang(lang string) string {
	l := strings.ToUpper(strings.TrimSpace(lang))
	if l == "" {
		return "ID"
	}
	return l
}

func translateTextPtrServiceGallery(text *string, fromLang, toLang string) *string {
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

func (u *serviceGalleryUsecaseImpl) uploadImage(ctx context.Context, image *multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	ext := filepath.Ext(image.Filename)
	objectName := fmt.Sprintf("images/service_gallery/%s%s", uuid.NewString(), ext)
	_, err = u.minioClient.PutObject(ctx, u.minioBucket, objectName, file, image.Size, minio.PutObjectOptions{ContentType: image.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (u *serviceGalleryUsecaseImpl) removeObject(ctx context.Context, objectPath string) error {
	objectName := normalizeObjectName(objectPath, u.publicBaseURL, u.minioBucket)
	if objectName == "" {
		return nil
	}
	return u.minioClient.RemoveObject(ctx, u.minioBucket, objectName, minio.RemoveObjectOptions{})
}
