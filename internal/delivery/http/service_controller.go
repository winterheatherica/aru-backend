package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ServiceController struct {
	Usecase usecase.ServiceUsecase
}

func NewServiceController(u usecase.ServiceUsecase) *ServiceController {
	return &ServiceController{
		Usecase: u,
	}
}

func (c *ServiceController) GetServicePage(ctx *fiber.Ctx) error {
	service := ctx.Query("service")
	lang := ctx.Query("lang", "id")

	if service == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "service is required",
		})
	}

	result, err := c.Usecase.GetServicePage(ctx.Context(), service, lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
