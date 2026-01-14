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

func (c *ArticleController) GetArticle(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "id")
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
