package http

import (
	"aru-backend/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type SiteContentController struct {
	Usecase usecase.SiteContentUsecase
}

func NewSiteContentController(u usecase.SiteContentUsecase) *SiteContentController {
	return &SiteContentController{
		Usecase: u,
	}
}

func (c *SiteContentController) GetSiteContent(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "ID")

	result, err := c.Usecase.GetSiteContent(ctx.Context(), lang)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(result)
}
