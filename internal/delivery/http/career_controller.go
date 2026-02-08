package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type CareerController struct {
	Usecase usecase.CareerUsecase
}

func NewCareerController(u usecase.CareerUsecase) *CareerController {
	return &CareerController{Usecase: u}
}

func (c *CareerController) GetCareers(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "ID")

	result, err := c.Usecase.GetCareers(ctx.Context(), lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
