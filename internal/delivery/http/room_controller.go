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
	lang := ctx.Query("lang", "ID")
	id := ctx.Params("id")

	result, err := c.Usecase.GetRoomByID(ctx.Context(), id, lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "room not found",
		})
	}

	return ctx.JSON(result)
}

func (c *RoomController) ResolveRoomID(ctx *fiber.Ctx) error {
	slug := ctx.Query("slug")

	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug is required",
		})
	}

	id, err := c.Usecase.ResolveRoomID(ctx.Context(), slug)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if id == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "room not found"})
	}

	return ctx.JSON(fiber.Map{"id": id})
}
