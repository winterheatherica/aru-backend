package config

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	bodyLimit := config.GetInt("web.bodyLimit")
	if bodyLimit <= 0 {
		bodyLimit = 50 * 1024 * 1024
	}

	var app = fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
		BodyLimit:    bodyLimit,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.GetString("web.cors.allowOrigins"),
		AllowMethods:     config.GetString("web.cors.allowMethods"),
		AllowHeaders:     config.GetString("web.cors.allowHeaders"),
		AllowCredentials: config.GetBool("web.cors.allowCredentials"),
	}))

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		return ctx.Status(code).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}
}
