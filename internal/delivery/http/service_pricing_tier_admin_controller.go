package http

import (
	"strings"

	"aru-backend/internal/model"
	adminusecase "aru-backend/internal/usecase/admin"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ServicePricingTierAdminController struct {
	Usecase adminusecase.ServicePricingTierUsecase
}

func NewServicePricingTierAdminController(u adminusecase.ServicePricingTierUsecase) *ServicePricingTierAdminController {
	return &ServicePricingTierAdminController{Usecase: u}
}

func (c *ServicePricingTierAdminController) ListByService(ctx *fiber.Ctx) error {
	service := strings.TrimSpace(ctx.Params("service"))
	lang := strings.ToUpper(strings.TrimSpace(ctx.Query("lang", "ID")))
	items, err := c.Usecase.ListByService(ctx.Context(), service, lang)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(items)
}

func (c *ServicePricingTierAdminController) GetByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	lang := strings.ToUpper(strings.TrimSpace(ctx.Query("lang", "ID")))
	item, err := c.Usecase.GetByID(ctx.Context(), id, lang)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.JSON(item)
}

func (c *ServicePricingTierAdminController) Create(ctx *fiber.Ctx) error {
	var input model.ServicePricingTierUpsertInput
	if err := ctx.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	created, err := c.Usecase.Create(ctx.Context(), input)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return ctx.Status(fiber.StatusCreated).JSON(created)
}

func (c *ServicePricingTierAdminController) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var input model.ServicePricingTierUpsertInput
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

func (c *ServicePricingTierAdminController) HardDelete(ctx *fiber.Ctx) error {
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
