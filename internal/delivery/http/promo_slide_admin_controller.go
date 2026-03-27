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

type PromoSlideAdminController struct {
	Usecase adminusecase.PromoSlideUsecase
}

func NewPromoSlideAdminController(u adminusecase.PromoSlideUsecase) *PromoSlideAdminController {
	return &PromoSlideAdminController{Usecase: u}
}

func (c *PromoSlideAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *PromoSlideAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *PromoSlideAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parsePromoSlideUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *PromoSlideAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	input, image, err := parsePromoSlideUpsertInput(ctx, false)
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

func (c *PromoSlideAdminController) HardDelete(ctx *fiber.Ctx) error {
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

func parsePromoSlideUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.PromoSlideUpsertInput, *multipart.FileHeader, error) {
	orderIndex, err := strconv.Atoi(ctx.FormValue("order_index", "0"))
	if err != nil {
		return model.PromoSlideUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "order_index must be a number")
	}

	activeRaw := strings.TrimSpace(ctx.FormValue("is_active", "true"))
	active := true
	if activeRaw != "" {
		v, perr := strconv.ParseBool(activeRaw)
		if perr != nil {
			return model.PromoSlideUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active must be boolean")
		}
		active = v
	}

	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))
	input := model.PromoSlideUpsertInput{
		Language:   language,
		Alt:        ptrOrNil(ctx.FormValue("alt")),
		Title:      ptrOrNil(ctx.FormValue("title")),
		OrderIndex: orderIndex,
		IsActive:   active,
	}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.PromoSlideUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}
	return input, image, nil
}
