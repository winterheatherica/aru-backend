package route

import (
	"aru-backend/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	AuthMiddleware fiber.Handler

	UserController *http.UserController

	HomeController        *http.HomeController
	AboutController       *http.AboutController
	ServiceController     *http.ServiceController
	ReservationController *http.ReservationController
	InformationController *http.InformationController
	CareerController      *http.CareerController

	RoomController *http.RoomController

	ArticleController  *http.ArticleController
	CategoryController *http.CategoryController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	c.App.Get("/api/home", c.HomeController.GetHome)
	c.App.Get("/api/about", c.AboutController.GetAbout)
	c.App.Get("/api/service", c.ServiceController.GetServicePageBatch)
	c.App.Get("/api/reservation", c.ReservationController.GetReservationPage)
	c.App.Get("/api/information", c.InformationController.GetInformation)
	c.App.Get("/api/career", c.CareerController.GetCareers)

	c.App.Get("/api/rooms/:slug", c.RoomController.GetRoomDetail)

	c.App.Get("/api/article/:id", c.ArticleController.GetArticle)
	c.App.Get("/api/category/:slug", c.CategoryController.GetCategory)

}

func (c *RouteConfig) SetupAuthRoute() {

}
