package http

import (
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"aru-backend/internal/model"
	adminusecase "aru-backend/internal/usecase/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NewsArticleAdminController struct {
	Usecase adminusecase.NewsArticleUsecase
}

func NewNewsArticleAdminController(u adminusecase.NewsArticleUsecase) *NewsArticleAdminController {
	return &NewsArticleAdminController{Usecase: u}
}

func (c *NewsArticleAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *NewsArticleAdminController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	item, err := c.Usecase.GetByID(ctx.Context(), id)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(item)
}

func (c *NewsArticleAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parseNewsArticleUpsertInput(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *NewsArticleAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	input, image, err := parseNewsArticleUpsertInput(ctx)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	updated, err := c.Usecase.Update(ctx.Context(), id, input, image)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(updated)
}

func (c *NewsArticleAdminController) HardDelete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := c.Usecase.HardDelete(ctx.Context(), id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.JSON(fiber.Map{"message": "deleted"})
}

func parseNewsArticleUpsertInput(ctx *fiber.Ctx) (model.NewsArticleUpsertInput, *multipart.FileHeader, error) {
	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))
	title := strings.TrimSpace(ctx.FormValue("title"))
	if title == "" {
		return model.NewsArticleUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "title is required")
	}

	content := strings.TrimSpace(ctx.FormValue("content"))
	if content == "" {
		return model.NewsArticleUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "content is required")
	}

	isActive := true
	if raw := strings.TrimSpace(ctx.FormValue("is_active", "true")); raw != "" {
		v, err := strconv.ParseBool(raw)
		if err != nil {
			return model.NewsArticleUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean")
		}
		isActive = v
	}

	publishedAt := time.Now().UTC()
	if raw := strings.TrimSpace(ctx.FormValue("published_at")); raw != "" {
		parsed, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return model.NewsArticleUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "published_at must be RFC3339")
		}
		publishedAt = parsed.UTC()
	}

	uploadedBy, err := parseOptionalUUID(ctx.FormValue("uploaded_by"))
	if err != nil {
		return model.NewsArticleUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "uploaded_by must be valid uuid")
	}

	input := model.NewsArticleUpsertInput{
		Language:        language,
		Title:           title,
		Slug:            strings.TrimSpace(ctx.FormValue("slug")),
		Content:         content,
		MetaTitle:       ptrOrNil(ctx.FormValue("meta_title")),
		MetaDescription: ptrOrNil(ctx.FormValue("meta_description")),
		MetaKeywords:    parseMetaKeywords(ctx.FormValue("meta_keywords")),
		CategoryIDs:     parseCategoryIDs(ctx.FormValue("category_ids")),
		UploadedBy:      uploadedBy,
		PublishedAt:     publishedAt,
		IsActive:        isActive,
	}

	image, _ := ctx.FormFile("image")
	return input, image, nil
}

func parseMetaKeywords(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}

func parseCategoryIDs(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		id := strings.TrimSpace(p)
		if id != "" {
			res = append(res, id)
		}
	}
	return res
}

func parseOptionalUUID(raw string) (*uuid.UUID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	id, err := uuid.Parse(trimmed)
	if err != nil {
		return nil, err
	}
	return &id, nil
}
