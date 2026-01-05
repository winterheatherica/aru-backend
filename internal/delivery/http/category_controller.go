package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	Usecase usecase.CategoryUsecase
}

func NewCategoryController(u usecase.CategoryUsecase) *CategoryController {
	return &CategoryController{
		Usecase: u,
	}
}

func (c *CategoryController) GetCategory(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")
	slug := ctx.Params("slug")

	result, err := c.Usecase.GetCategoryBySlug(ctx.Context(), slug, lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "category not found",
		})
	}

	return ctx.JSON(result)
}
