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

type HeroAdminController struct {
	Usecase adminusecase.HeroUsecase
}

func NewHeroAdminController(u adminusecase.HeroUsecase) *HeroAdminController {
	return &HeroAdminController{Usecase: u}
}

func (c *HeroAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *HeroAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *HeroAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parseHeroUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *HeroAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	input, image, err := parseHeroUpsertInput(ctx, false)
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

func (c *HeroAdminController) HardDelete(ctx *fiber.Ctx) error {
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

func parseHeroUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.HeroUpsertInput, *multipart.FileHeader, error) {
	orderIndex, err := strconv.Atoi(ctx.FormValue("order_index", "0"))
	if err != nil {
		return model.HeroUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "order_index must be a number")
	}

	isActiveRaw := strings.TrimSpace(ctx.FormValue("is_active", "true"))
	isActive := true
	if isActiveRaw != "" {
		v, parseErr := strconv.ParseBool(isActiveRaw)
		if parseErr != nil {
			return model.HeroUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean")
		}
		isActive = v
	}

	banner := strings.TrimSpace(ctx.FormValue("banner", "POLISH"))
	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))

	descriptionRaw := ctx.FormValue("description")
	if strings.TrimSpace(descriptionRaw) == "" {
		descriptionRaw = ctx.FormValue("cta_label")
	}

	input := model.HeroUpsertInput{
		Language:    language,
		Alt:         ptrOrNil(ctx.FormValue("alt")),
		Title:       ptrOrNil(ctx.FormValue("title")),
		Description: ptrOrNil(descriptionRaw),
		Banner:      banner,
		OrderIndex:  orderIndex,
		IsActive:    isActive,
	}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.HeroUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}

	return input, image, nil
}

func ptrOrNil(s string) *string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
