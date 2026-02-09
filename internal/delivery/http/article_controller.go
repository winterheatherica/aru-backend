package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ArticleController struct {
	Usecase usecase.ArticleUsecase
}

func NewArticleController(u usecase.ArticleUsecase) *ArticleController {
	return &ArticleController{
		Usecase: u,
	}
}

func (c *ArticleController) GetArticleDetail(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "ID")
	id := ctx.Params("id")

	result, err := c.Usecase.GetArticleByID(ctx.Context(), id, lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "article not found",
		})
	}

	return ctx.JSON(result)
}

func (c *ArticleController) ResolveArticleID(ctx *fiber.Ctx) error {
	slug := ctx.Query("slug")

	if slug == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug is required",
		})
	}

	id, err := c.Usecase.ResolveArticleID(ctx.Context(), slug)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if id == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "article not found"})
	}

	return ctx.JSON(fiber.Map{"id": id})
}
