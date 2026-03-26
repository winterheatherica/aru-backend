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

type ClientAdminController struct {
	Usecase adminusecase.ClientUsecase
}

func NewClientAdminController(u adminusecase.ClientUsecase) *ClientAdminController {
	return &ClientAdminController{Usecase: u}
}

func (c *ClientAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *ClientAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *ClientAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parseClientUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *ClientAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	input, image, err := parseClientUpsertInput(ctx, false)
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

func (c *ClientAdminController) HardDelete(ctx *fiber.Ctx) error {
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

func parseClientUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.ClientUpsertInput, *multipart.FileHeader, error) {
	orderIndex, err := strconv.Atoi(ctx.FormValue("order_index", "0"))
	if err != nil {
		return model.ClientUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "order_index must be a number")
	}

	activeRaw := strings.TrimSpace(ctx.FormValue("is_active_client_scroller", "true"))
	active := true
	if activeRaw != "" {
		v, perr := strconv.ParseBool(activeRaw)
		if perr != nil {
			return model.ClientUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active_client_scroller must be boolean")
		}
		active = v
	}

	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))
	input := model.ClientUpsertInput{Language: language, Alt: ptrOrNil(ctx.FormValue("alt")), Title: ptrOrNil(ctx.FormValue("title")), Description: ptrOrNil(ctx.FormValue("description")), OrderIndex: orderIndex, IsActiveClientScroller: active}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.ClientUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}
	return input, image, nil
}
