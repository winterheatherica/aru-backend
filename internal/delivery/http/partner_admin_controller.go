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

type PartnerAdminController struct {
	Usecase adminusecase.PartnerUsecase
}

func NewPartnerAdminController(u adminusecase.PartnerUsecase) *PartnerAdminController {
	return &PartnerAdminController{Usecase: u}
}

func (c *PartnerAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *PartnerAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *PartnerAdminController) Create(ctx *fiber.Ctx) error {
	input, image, err := parsePartnerUpsertInput(ctx, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input, image)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *PartnerAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	input, image, err := parsePartnerUpsertInput(ctx, false)
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

func (c *PartnerAdminController) HardDelete(ctx *fiber.Ctx) error {
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

func parsePartnerUpsertInput(ctx *fiber.Ctx, imageRequired bool) (model.PartnerUpsertInput, *multipart.FileHeader, error) {
	orderIndex, err := strconv.Atoi(ctx.FormValue("order_index", "0"))
	if err != nil {
		return model.PartnerUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "order_index must be a number")
	}

	gridRaw := strings.TrimSpace(ctx.FormValue("is_active_partner_grid", "true"))
	grid := true
	if gridRaw != "" {
		v, perr := strconv.ParseBool(gridRaw)
		if perr != nil {
			return model.PartnerUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active_partner_grid must be boolean")
		}
		grid = v
	}

	scrollerRaw := strings.TrimSpace(ctx.FormValue("is_active_partner_scroller", "true"))
	scroller := true
	if scrollerRaw != "" {
		v, perr := strconv.ParseBool(scrollerRaw)
		if perr != nil {
			return model.PartnerUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "is_active_partner_scroller must be boolean")
		}
		scroller = v
	}

	language := strings.ToUpper(strings.TrimSpace(ctx.FormValue("language", "ID")))
	input := model.PartnerUpsertInput{
		Language:                language,
		Alt:                     ptrOrNil(ctx.FormValue("alt")),
		Title:                   ptrOrNil(ctx.FormValue("title")),
		Description:             ptrOrNil(ctx.FormValue("description")),
		OrderIndex:              orderIndex,
		IsActivePartnerGrid:     grid,
		IsActivePartnerScroller: scroller,
	}

	image, _ := ctx.FormFile("image")
	if imageRequired && image == nil {
		return model.PartnerUpsertInput{}, nil, fiber.NewError(fiber.StatusBadRequest, "image is required")
	}
	return input, image, nil
}
