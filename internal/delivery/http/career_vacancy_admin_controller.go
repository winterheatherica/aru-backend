package http

import (
	"strings"

	"aru-backend/internal/model"
	adminusecase "aru-backend/internal/usecase/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CareerVacancyAdminController struct {
	Usecase adminusecase.CareerVacancyUsecase
}

func NewCareerVacancyAdminController(u adminusecase.CareerVacancyUsecase) *CareerVacancyAdminController {
	return &CareerVacancyAdminController{Usecase: u}
}

func (c *CareerVacancyAdminController) List(ctx *fiber.Ctx) error {
	items, err := c.Usecase.List(ctx.Context())
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *CareerVacancyAdminController) GetByID(ctx *fiber.Ctx) error {
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

func (c *CareerVacancyAdminController) Create(ctx *fiber.Ctx) error {
	var input model.CareerVacancyUpsertInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *CareerVacancyAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var input model.CareerVacancyUpsertInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	updated, err := c.Usecase.Update(ctx.Context(), id, input)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(updated)
}

func (c *CareerVacancyAdminController) HardDelete(ctx *fiber.Ctx) error {
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
