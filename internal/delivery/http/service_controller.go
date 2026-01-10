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

func (c *ServiceController) GetServicePageBatch(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")

	services := []string{
		"ARUCONTRACTOR",
		"ARUDIGITAL",
		"ARUHEALTHCARE",
		"ARULOG",
		"ARUSOLUTION",
		"ARUSOURCE",
		"ARUSPACE",
		"ARUTRANS",
	}

	result := make(map[string]interface{})

	for _, service := range services {
		data, err := c.Usecase.GetServicePage(ctx.Context(), service, lang)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"service": service,
			})
		}

		result[service] = data
	}

	return ctx.JSON(result)
}
