package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type RoomController struct {
	Usecase usecase.RoomUsecase
}

func NewRoomController(u usecase.RoomUsecase) *RoomController {
	return &RoomController{
		Usecase: u,
	}
}

func (c *RoomController) GetRoomDetail(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")
	slug := ctx.Params("slug")

	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug is required",
		})
	}

	result, err := c.Usecase.GetRoomDetail(ctx.Context(), slug, lang)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
