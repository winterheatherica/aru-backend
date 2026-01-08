package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type InformationController struct {
	Usecase usecase.InformationUsecase
}

func NewInformationController(u usecase.InformationUsecase) *InformationController {
	return &InformationController{
		Usecase: u,
	}
}

func (c *InformationController) GetInformation(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")
	page := ctx.QueryInt("page", 1)

	var year *int
	if y := ctx.QueryInt("year", 0); y != 0 {
		year = &y
	}

	result, err := c.Usecase.GetInformation(
		ctx.Context(),
		lang,
		year,
		page,
	)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
