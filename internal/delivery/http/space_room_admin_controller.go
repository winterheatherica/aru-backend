package http

import (
	"mime/multipart"
	"strconv"
	"strings"

	"aru-backend/internal/model"
	adminusecase "aru-backend/internal/usecase/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SpaceRoomAdminController struct{ Usecase adminusecase.SpaceRoomUsecase }

func NewSpaceRoomAdminController(u adminusecase.SpaceRoomUsecase) *SpaceRoomAdminController {
	return &SpaceRoomAdminController{Usecase: u}
}

func (c *SpaceRoomAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil { return fiber.NewError(fiber.StatusInternalServerError, err.Error()) }
	return ctx.JSON(items)
}

func (c *SpaceRoomAdminController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id")); if err != nil { return fiber.NewError(fiber.StatusBadRequest, "invalid id") }
	item, err := c.Usecase.GetByID(ctx.Context(), id)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return ctx.JSON(item)
}

func (c *SpaceRoomAdminController) Create(ctx *fiber.Ctx) error {
	input, files, err := parseSpaceRoomInput(ctx)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	item, err := c.Usecase.Create(ctx.Context(), input, files)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return ctx.Status(fiber.StatusCreated).JSON(item)
}

func (c *SpaceRoomAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id")); if err != nil { return fiber.NewError(fiber.StatusBadRequest, "invalid id") }
	input, files, err := parseSpaceRoomInput(ctx)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	item, err := c.Usecase.Update(ctx.Context(), id, input, files)
	if err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return ctx.JSON(item)
}

func (c *SpaceRoomAdminController) HardDelete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id")); if err != nil { return fiber.NewError(fiber.StatusBadRequest, "invalid id") }
	if err := c.Usecase.HardDelete(ctx.Context(), id); err != nil { return fiber.NewError(fiber.StatusBadRequest, err.Error()) }
	return ctx.JSON(fiber.Map{"message":"deleted"})
}

func parseSpaceRoomInput(ctx *fiber.Ctx) (model.SpaceRoomAdminUpsertInput, []*multipart.FileHeader, error) {
	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))
	code := strings.TrimSpace(ctx.FormValue("code"))
	title := strings.TrimSpace(ctx.FormValue("title"))
	description := strings.TrimSpace(ctx.FormValue("description"))
	if code == "" || title == "" { return model.SpaceRoomAdminUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "code and title are required") }

	isActive := true
	if v := strings.TrimSpace(ctx.FormValue("is_active", "true")); v != "" {
		b, err := strconv.ParseBool(v); if err != nil { return model.SpaceRoomAdminUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean") }
		isActive = b
	}

	var capacity *int
	if raw := strings.TrimSpace(ctx.FormValue("capacity")); raw != "" {
		val, err := strconv.Atoi(raw); if err != nil { return model.SpaceRoomAdminUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "capacity must be number") }
		capacity = &val
	}
	var floor *int
	if raw := strings.TrimSpace(ctx.FormValue("floor")); raw != "" {
		val, err := strconv.Atoi(raw); if err != nil { return model.SpaceRoomAdminUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "floor must be number") }
		floor = &val
	}

	var uploadedBy *uuid.UUID
	if raw := strings.TrimSpace(ctx.FormValue("uploaded_by")); raw != "" {
		id, err := uuid.Parse(raw); if err != nil { return model.SpaceRoomAdminUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "uploaded_by must be valid uuid") }
		uploadedBy = &id
	}

	thumbnailIndex := 0
	if raw := strings.TrimSpace(ctx.FormValue("thumbnail_index", "0")); raw != "" {
		idx, err := strconv.Atoi(raw); if err == nil { thumbnailIndex = idx }
	}

	metaKeywords := splitComma(ctx.FormValue("meta_keywords"))
	facilities := splitComma(ctx.FormValue("facilities"))

	form, _ := ctx.MultipartForm()
	files := []*multipart.FileHeader{}
	if form != nil {
		if list, ok := form.File["images"]; ok { files = list }
		if len(files) == 0 {
			for _, group := range form.File {
				if len(group) > 0 { files = append(files, group...) }
			}
		}
	}

	alts := splitPipe(ctx.FormValue("image_alts"))
	titles := splitPipe(ctx.FormValue("image_titles"))
	imgInputs := make([]model.SpaceRoomAdminImageInput, 0)
	max := len(files)
	for i := 0; i < max; i++ {
		alt := ""
		title := ""
		if i < len(alts) { alt = alts[i] }
		if i < len(titles) { title = titles[i] }
		imgInputs = append(imgInputs, model.SpaceRoomAdminImageInput{Alt: alt, Title: title})
	}

	return model.SpaceRoomAdminUpsertInput{
		Language: language, Code: code, Title: title, Description: description,
		Facilities: facilities, MetaKeywords: metaKeywords,
		Capacity: capacity, Floor: floor, Address: ptrOrNil(ctx.FormValue("address")),
		IsActive: isActive, UploadedBy: uploadedBy,
		Images: imgInputs, ThumbnailIndex: thumbnailIndex,
	}, files, nil
}

func splitComma(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" { return nil }
	parts := strings.Split(trimmed, ",")
	res := make([]string,0,len(parts))
	for _, p := range parts { t := strings.TrimSpace(p); if t != "" { res = append(res, t) } }
	return res
}

func splitPipe(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" { return nil }
	parts := strings.Split(trimmed, "||")
	res := make([]string,0,len(parts))
	for _, p := range parts { res = append(res, strings.TrimSpace(p)) }
	return res
}
