package admin

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strings"

	"aru-backend/internal/entity"
	"aru-backend/internal/model"
	"aru-backend/internal/model/converter"
	"aru-backend/internal/repository"

	"github.com/bregydoc/gtranslate"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/minio/minio-go/v7"
)

type SpaceRoomUsecase interface {
	List(ctx context.Context) ([]model.SpaceRoomAdminItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.SpaceRoomAdminItem, error)
	Create(ctx context.Context, input model.SpaceRoomAdminUpsertInput, images []*multipart.FileHeader) (*model.SpaceRoomAdminItem, error)
	Update(ctx context.Context, id uuid.UUID, input model.SpaceRoomAdminUpsertInput, images []*multipart.FileHeader) (*model.SpaceRoomAdminItem, error)
	HardDelete(ctx context.Context, id uuid.UUID) error
}

type spaceRoomUsecaseImpl struct {
	repo          repository.SpaceRoomRepository
	minioClient   *minio.Client
	minioBucket   string
	publicBaseURL string
}

func NewSpaceRoomUsecase(repo repository.SpaceRoomRepository, minioClient *minio.Client, minioBucket, publicBaseURL string) SpaceRoomUsecase {
	return &spaceRoomUsecaseImpl{repo: repo, minioClient: minioClient, minioBucket: minioBucket, publicBaseURL: strings.TrimSuffix(publicBaseURL, "/")}
}

func (u *spaceRoomUsecaseImpl) List(ctx context.Context) ([]model.SpaceRoomAdminItem, error) {
	rooms, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]model.SpaceRoomAdminItem, 0, len(rooms))
	for _, r := range rooms {
		item := u.toAdminItem(r)
		u.attachUploaderName(ctx, &item)
		res = append(res, item)
	}
	return res, nil
}

func (u *spaceRoomUsecaseImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.SpaceRoomAdminItem, error) {
	room, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item := u.toAdminItem(*room)
	u.attachUploaderName(ctx, &item)
	return &item, nil
}

func (u *spaceRoomUsecaseImpl) Create(ctx context.Context, input model.SpaceRoomAdminUpsertInput, images []*multipart.FileHeader) (*model.SpaceRoomAdminItem, error) {
	lang := normalizeLangSpaceRoom(input.Language)
	title := strings.TrimSpace(input.Title)
	slug := buildSlugSpaceRoom(title)
	slug, _ = u.ensureUniqueSlug(ctx, slug, lang, nil)

	if strings.TrimSpace(input.Code) == "" {
		return nil, fmt.Errorf("code is required")
	}
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("at least one image is required")
	}

	thumb := input.ThumbnailIndex
	if thumb < 0 || thumb >= len(images) {
		thumb = 0
	}

	idDesc := strings.TrimSpace(input.Description)
	idMetaTitle := deriveMetaTitleSpaceRoom(title)
	idMetaDesc := deriveMetaDescriptionSpaceRoom(idDesc)

	translation := entity.SpaceRoomTranslation{
		Language:        lang,
		Title:           title,
		Description:     stringPtrOrNil(idDesc),
		Facilities:      toPQStrings(input.Facilities),
		Slug:            slug,
		MetaTitle:       idMetaTitle,
		MetaDescription: idMetaDesc,
		MetaKeywords:    toPQStrings(input.MetaKeywords),
	}

	enTitle := translateTextSpaceRoom(title, "id", "en")
	enSlug := buildSlugSpaceRoom(enTitle)
	if enSlug == slug {
		enSlug = enSlug + "-en"
	}
	enSlug, _ = u.ensureUniqueSlug(ctx, enSlug, "EN", nil)
	enDesc := translateTextSpaceRoom(idDesc, "id", "en")

	translations := []entity.SpaceRoomTranslation{translation}
	if lang == "ID" {
		translations = append(translations, entity.SpaceRoomTranslation{
			Language:        "EN",
			Title:           enTitle,
			Description:     stringPtrOrNil(enDesc),
			Facilities:      toPQStrings(translateListSpaceRoom(input.Facilities, "id", "en")),
			Slug:            enSlug,
			MetaTitle:       deriveMetaTitleSpaceRoom(enTitle),
			MetaDescription: deriveMetaDescriptionSpaceRoom(enDesc),
			MetaKeywords:    toPQStrings(translateListSpaceRoom(input.MetaKeywords, "id", "en")),
		})
	}

	room := &entity.SpaceRoom{
		Code:         strings.TrimSpace(input.Code),
		Capacity:     input.Capacity,
		Floor:        input.Floor,
		Address:      input.Address,
		UploadedBy:   input.UploadedBy,
		IsActive:     input.IsActive,
		Translations: translations,
	}

	if err := u.repo.Create(ctx, room); err != nil {
		return nil, err
	}

	if err := u.saveImages(ctx, room.ID, input, images, thumb, room.UploadedBy, lang); err != nil {
		return nil, err
	}

	fresh, err := u.repo.FindByID(ctx, room.ID)
	if err != nil {
		item := u.toAdminItem(*room)
		u.attachUploaderName(ctx, &item)
		return &item, nil
	}
	item := u.toAdminItem(*fresh)
	u.attachUploaderName(ctx, &item)
	return &item, nil
}

func (u *spaceRoomUsecaseImpl) Update(ctx context.Context, id uuid.UUID, input model.SpaceRoomAdminUpsertInput, images []*multipart.FileHeader) (*model.SpaceRoomAdminItem, error) {
	room, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	lang := normalizeLangSpaceRoom(input.Language)
	title := strings.TrimSpace(input.Title)
	slug := buildSlugSpaceRoom(title)
	slug, _ = u.ensureUniqueSlug(ctx, slug, lang, &id)

	room.Code = strings.TrimSpace(input.Code)
	room.Capacity = input.Capacity
	room.Floor = input.Floor
	room.Address = input.Address
	room.IsActive = input.IsActive
	if input.UploadedBy != nil {
		room.UploadedBy = input.UploadedBy
	}
	if err := u.repo.Update(ctx, room); err != nil {
		return nil, err
	}

	tr := &entity.SpaceRoomTranslation{
		RoomID:           room.ID,
		Language:         lang,
		Title:            title,
		Description:      stringPtrOrNil(input.Description),
		Facilities:       toPQStrings(input.Facilities),
		Slug:             slug,
		MetaTitle:        deriveMetaTitleSpaceRoom(title),
		MetaDescription:  deriveMetaDescriptionSpaceRoom(input.Description),
		MetaKeywords:     toPQStrings(input.MetaKeywords),
	}
	if err := u.repo.UpsertTranslation(ctx, tr); err != nil {
		return nil, err
	}

	if len(images) > 0 {
		thumb := input.ThumbnailIndex
		if thumb < 0 || thumb >= len(images) {
			thumb = 0
		}
		if err := u.repo.DeleteImagesByRoomID(ctx, room.ID); err != nil {
			return nil, err
		}
		if err := u.saveImages(ctx, room.ID, input, images, thumb, room.UploadedBy, lang); err != nil {
			return nil, err
		}
	}

	fresh, err := u.repo.FindByID(ctx, room.ID)
	if err != nil {
		item := u.toAdminItem(*room)
		u.attachUploaderName(ctx, &item)
		return &item, nil
	}
	item := u.toAdminItem(*fresh)
	u.attachUploaderName(ctx, &item)
	return &item, nil
}

func (u *spaceRoomUsecaseImpl) HardDelete(ctx context.Context, id uuid.UUID) error { return u.repo.DeleteByID(ctx, id) }

func (u *spaceRoomUsecaseImpl) saveImages(ctx context.Context, roomID uuid.UUID, input model.SpaceRoomAdminUpsertInput, files []*multipart.FileHeader, thumb int, uploadedBy *uuid.UUID, lang string) error {
	for i, f := range files {
		objectName, err := u.uploadImage(ctx, f)
		if err != nil {
			return err
		}
		img := &entity.SpaceRoomImage{RoomID: roomID, UploadedBy: uploadedBy, ImageURL: objectName, IsActive: true, OrderIndex: i + 1}
		if err := u.repo.CreateImage(ctx, img); err != nil {
			return err
		}
		if i == thumb {
			if r, e := u.repo.FindByID(ctx, roomID); e == nil {
				r.MainImageURL = &objectName
				_ = u.repo.Update(ctx, r)
			}
		}

		alt := ""
		ttl := ""
		if i < len(input.Images) {
			alt = strings.TrimSpace(input.Images[i].Alt)
			ttl = strings.TrimSpace(input.Images[i].Title)
		}
		if alt == "" {
			alt = fmt.Sprintf("%s image %d", strings.TrimSpace(input.Title), i+1)
		}
		if ttl == "" {
			ttl = fmt.Sprintf("%s image %d", strings.TrimSpace(input.Title), i+1)
		}
		_ = u.repo.UpsertImageTranslation(ctx, &entity.SpaceRoomImageTranslation{ImageID: img.ID, Language: lang, Alt: &alt, Title: &ttl})
		if strings.EqualFold(lang, "ID") {
			enAlt := translateTextSpaceRoom(alt, "id", "en")
			enTitle := translateTextSpaceRoom(ttl, "id", "en")
			_ = u.repo.UpsertImageTranslation(ctx, &entity.SpaceRoomImageTranslation{ImageID: img.ID, Language: "EN", Alt: &enAlt, Title: &enTitle})
		}
	}
	return nil
}

func (u *spaceRoomUsecaseImpl) toAdminItem(room entity.SpaceRoom) model.SpaceRoomAdminItem {
	trs := make([]model.SpaceRoomAdminTranslation, 0, len(room.Translations))
	for _, tr := range room.Translations {
		trs = append(trs, model.SpaceRoomAdminTranslation{Language: tr.Language, Slug: tr.Slug, Title: tr.Title, Description: tr.Description, Facilities: []string(tr.Facilities), MetaTitle: tr.MetaTitle, MetaDescription: tr.MetaDescription, MetaKeywords: []string(tr.MetaKeywords)})
	}
	imgs := make([]model.SpaceRoomAdminImage, 0, len(room.Images))
	for _, img := range room.Images {
		var alt, title *string
		if len(img.Translations) > 0 {
			alt = img.Translations[0].Alt
			title = img.Translations[0].Title
		}
		isThumb := room.MainImageURL != nil && *room.MainImageURL == img.ImageURL
		url := converter.BuildAssetURL(u.publicBaseURL, img.ImageURL)
		imgs = append(imgs, model.SpaceRoomAdminImage{ID: img.ID, ImageURL: url, OrderIndex: img.OrderIndex, IsActive: img.IsActive, IsThumbnail: isThumb, Alt: alt, Title: title})
	}
	var mainURL *string
	if room.MainImageURL != nil {
		u := converter.BuildAssetURL(u.publicBaseURL, *room.MainImageURL)
		mainURL = &u
	}
	return model.SpaceRoomAdminItem{ID: room.ID, Code: room.Code, Capacity: room.Capacity, Floor: room.Floor, Address: room.Address, MainImageURL: mainURL, UploadedBy: room.UploadedBy, IsActive: room.IsActive, CreatedAt: room.CreatedAt, UpdatedAt: room.UpdatedAt, Translations: trs, Images: imgs}
}

func (u *spaceRoomUsecaseImpl) ensureUniqueSlug(ctx context.Context, slug, lang string, excludeID *uuid.UUID) (string, error) {
	base := strings.TrimSpace(slug)
	if base == "" {
		base = "room"
	}
	candidate := base
	for i := 1; i <= 200; i++ {
		exists, err := u.repo.IsSlugExists(ctx, candidate, lang, excludeID)
		if err != nil {
			return candidate, err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
	return candidate, nil
}

func normalizeLangSpaceRoom(lang string) string {
	v := strings.ToUpper(strings.TrimSpace(lang))
	if v == "" {
		return "ID"
	}
	return v
}

var nonAlnumSpaceRoom = regexp.MustCompile(`[^a-z0-9]+`)

func buildSlugSpaceRoom(text string) string {
	t := strings.ToLower(strings.TrimSpace(text))
	t = nonAlnumSpaceRoom.ReplaceAllString(t, "-")
	t = strings.Trim(t, "-")
	if t == "" {
		return "space-room"
	}
	return t
}

func deriveMetaTitleSpaceRoom(title string) *string {
	t := strings.TrimSpace(title)
	if t == "" {
		return nil
	}
	r := []rune(t)
	if len(r) > 70 {
		t = strings.TrimSpace(string(r[:70]))
	}
	return &t
}

var htmlStripSpaceRoom = regexp.MustCompile(`<[^>]+>`)

func deriveMetaDescriptionSpaceRoom(content string) *string {
	t := strings.TrimSpace(content)
	if t == "" {
		return nil
	}
	t = htmlStripSpaceRoom.ReplaceAllString(t, " ")
	t = strings.Join(strings.Fields(t), " ")
	r := []rune(t)
	if len(r) > 160 {
		t = strings.TrimSpace(string(r[:160]))
	}
	if t == "" {
		return nil
	}
	return &t
}

func toPQStrings(values []string) pq.StringArray {
	res := make([]string, 0, len(values))
	for _, v := range values {
		t := strings.TrimSpace(v)
		if t != "" {
			res = append(res, t)
		}
	}
	if len(res) == 0 {
		return nil
	}
	return pq.StringArray(res)
}

func translateTextSpaceRoom(text, from, to string) string {
	t := strings.TrimSpace(text)
	if t == "" {
		return ""
	}
	res, err := gtranslate.TranslateWithParams(t, gtranslate.TranslationParams{From: from, To: to})
	if err != nil || strings.TrimSpace(res) == "" {
		return t
	}
	return strings.TrimSpace(res)
}

func translateListSpaceRoom(values []string, from, to string) []string {
	res := make([]string, 0, len(values))
	for _, v := range values {
		t := translateTextSpaceRoom(v, from, to)
		if strings.TrimSpace(t) != "" {
			res = append(res, t)
		}
	}
	return res
}

func (u *spaceRoomUsecaseImpl) uploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()
	ext := filepath.Ext(file.Filename)
	object := fmt.Sprintf("images/space_rooms/%s%s", uuid.NewString(), ext)
	_, err = u.minioClient.PutObject(ctx, u.minioBucket, object, f, file.Size, minio.PutObjectOptions{ContentType: file.Header.Get("Content-Type")})
	if err != nil {
		return "", err
	}
	return object, nil
}

func stringPtrOrNil(v string) *string {
	t := strings.TrimSpace(v)
	if t == "" {
		return nil
	}
	return &t
}

func (u *spaceRoomUsecaseImpl) attachUploaderName(ctx context.Context, item *model.SpaceRoomAdminItem) {
	if item == nil || item.UploadedBy == nil {
		return
	}
	p, err := u.repo.FindUploaderByRoomID(ctx, item.ID)
	if err != nil || p == nil {
		return
	}
	name := ""
	if p.FullName != nil {
		name = strings.TrimSpace(*p.FullName)
	}
	if name == "" && p.Username != nil {
		name = strings.TrimSpace(*p.Username)
	}
	if name == "" {
		name = strings.TrimSpace(p.Email)
	}
	if name != "" {
		item.UploadedByName = &name
	}
}

