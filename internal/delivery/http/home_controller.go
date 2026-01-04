package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type HomeController struct {
	Usecase usecase.HomeUsecase
}

func NewHomeController(u usecase.HomeUsecase) *HomeController {
	return &HomeController{
		Usecase: u,
	}
}

func (c *HomeController) GetHome(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")

	result, err := c.Usecase.GetHome(ctx.Context(), lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
