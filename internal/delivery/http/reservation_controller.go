package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ReservationController struct {
	Usecase usecase.ReservationUsecase
}

func NewReservationController(u usecase.ReservationUsecase) *ReservationController {
	return &ReservationController{Usecase: u}
}

func (c *ReservationController) GetReservationPage(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")

	result, err := c.Usecase.GetReservationPage(ctx.Context(), lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
