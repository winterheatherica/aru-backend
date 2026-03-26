package http

import (
	"aru-backend/internal/usecase"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

func NewUserController(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req loginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	identifier := req.Identifier
	if identifier == "" {
		if req.Email != "" {
			identifier = req.Email
		} else {
			identifier = req.Username
		}
	}

	result, err := c.UseCase.LoginWithPassword(
		identifier,
		req.Password,
		ctx.IP(),
		ctx.Get("User-Agent"),
	)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    result.AccessToken,
		HTTPOnly: true,
		Secure:   false,
		Path:     "/",
		MaxAge:   int((24 * time.Hour).Seconds()),
		SameSite: "Lax",
	})

	return ctx.JSON(result)
}

func (c *UserController) Me(ctx *fiber.Ctx) error {
	token := ctx.Cookies("session")
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := c.UseCase.GetCurrentUserByAccessToken(token)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	return ctx.JSON(user)
}

func (c *UserController) Logout(ctx *fiber.Ctx) error {
	token := ctx.Cookies("session")
	_ = c.UseCase.LogoutByAccessToken(token)

	ctx.Cookie(&fiber.Cookie{Name: "session", Value: "", Path: "/", MaxAge: -1})
	return ctx.JSON(fiber.Map{"message": "ok"})
}
