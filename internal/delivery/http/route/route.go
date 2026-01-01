package route

import (
	"aru-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler

	UserController        *http.UserController
	SiteContentController *http.SiteContentController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Get("/api/site", c.SiteContentController.GetSiteContent)
}
func (c *RouteConfig) SetupAuthRoute() {

}
