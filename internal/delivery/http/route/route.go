package route

import (
	"aru-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler

	UserController   *http.UserController
	HomeController   *http.HomeController
	AboutController  *http.AboutController
	CareerController *http.CareerController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/api/home", c.HomeController.GetHome)
	c.App.Get("/api/about", c.AboutController.GetAbout)
	c.App.Get("/api/career", c.CareerController.GetCareers)
}

func (c *RouteConfig) SetupAuthRoute() {

}
