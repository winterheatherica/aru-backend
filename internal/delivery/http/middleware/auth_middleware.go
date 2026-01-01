package middleware

import "github.com/gofiber/fiber/v2"

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

type AuthMiddleware struct{}

func (m *AuthMiddleware) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return ctx.Next()
	}
}
