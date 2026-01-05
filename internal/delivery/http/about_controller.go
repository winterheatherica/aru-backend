package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AboutController struct {
	Usecase usecase.AboutUsecase
}

func NewAboutController(u usecase.AboutUsecase) *AboutController {
	return &AboutController{
		Usecase: u,
	}
}

func (c *AboutController) GetAbout(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")

	result, err := c.Usecase.GetAbout(ctx.Context(), lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
