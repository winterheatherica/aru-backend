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

type AwardAdminController struct {
	Usecase adminusecase.AwardUsecase
}

func NewAwardAdminController(u adminusecase.AwardUsecase) *AwardAdminController {
	return &AwardAdminController{Usecase: u}
}

func (c *AwardAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *AwardAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *AwardAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parseAwardUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *AwardAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	input, image, err := parseAwardUpsertInput(ctx, false)
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

func (c *AwardAdminController) HardDelete(ctx *fiber.Ctx) error {
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

func parseAwardUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.AwardUpsertInput, *multipart.FileHeader, error) {
	orderIndex, err := strconv.Atoi(ctx.FormValue("order_index", "0"))
	if err != nil {
		return model.AwardUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "order_index must be a number")
	}

	yearRaw := strings.TrimSpace(ctx.FormValue("year"))
	var yearPtr *int
	if yearRaw != "" {
		y, parseYearErr := strconv.Atoi(yearRaw)
		if parseYearErr != nil {
			return model.AwardUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "year must be a number")
		}
		yearPtr = &y
	}

	isActiveRaw := strings.TrimSpace(ctx.FormValue("is_active", "true"))
	isActive := true
	if isActiveRaw != "" {
		v, parseErr := strconv.ParseBool(isActiveRaw)
		if parseErr != nil {
			return model.AwardUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean")
		}
		isActive = v
	}

	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))

	input := model.AwardUpsertInput{
		Language:    language,
		Alt:         ptrOrNil(ctx.FormValue("alt")),
		Title:       ptrOrNil(ctx.FormValue("title")),
		Label:       ptrOrNil(ctx.FormValue("label")),
		Description: ptrOrNil(ctx.FormValue("description")),
		Year:        yearPtr,
		OrderIndex:  orderIndex,
		IsActive:    isActive,
	}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.AwardUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}

	return input, image, nil
}
