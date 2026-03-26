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

	HeroAdminController          *http.HeroAdminController
	AwardAdminController         *http.AwardAdminController
	CareerVacancyAdminController *http.CareerVacancyAdminController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/auth/login", c.UserController.Login)
	c.App.Get("/api/me", c.UserController.Me)
	c.App.Post("/api/auth/logout", c.UserController.Logout)

	c.App.Get("/api/home", c.HomeController.GetHome)
	c.App.Get("/api/about", c.AboutController.GetAbout)
	c.App.Get("/api/service", c.ServiceController.GetServicePageBatch)
	c.App.Get("/api/reservation", c.ReservationController.GetReservationPage)
	c.App.Get("/api/information", c.InformationController.GetInformation)
	c.App.Get("/api/career", c.CareerController.GetCareers)

	c.App.Get("/api/room/resolve", c.RoomController.ResolveRoomID)
	c.App.Get("/api/room/:id", c.RoomController.GetRoomDetail)

	c.App.Get("/api/article/resolve", c.ArticleController.ResolveArticleID)
	c.App.Get("/api/article/:id", c.ArticleController.GetArticleDetail)

	c.App.Get("/api/category/:slug", c.CategoryController.GetCategory)
}

func (c *RouteConfig) SetupAuthRoute() {
	admin := c.App.Group("/api/admin")
	if c.AuthMiddleware != nil {
		admin.Use(c.AuthMiddleware)
	}

	admin.Get("/hero", c.HeroAdminController.List)
	admin.Get("/hero/:id", c.HeroAdminController.GetByID)
	admin.Post("/hero", c.HeroAdminController.Create)
	admin.Put("/hero/:id", c.HeroAdminController.Update)
	admin.Delete("/hero/:id", c.HeroAdminController.HardDelete)

	admin.Get("/award", c.AwardAdminController.List)
	admin.Get("/award/:id", c.AwardAdminController.GetByID)
	admin.Post("/award", c.AwardAdminController.Create)
	admin.Put("/award/:id", c.AwardAdminController.Update)
	admin.Delete("/award/:id", c.AwardAdminController.HardDelete)

	admin.Get("/career-vacancy", c.CareerVacancyAdminController.List)
	admin.Get("/career-vacancy/:id", c.CareerVacancyAdminController.GetByID)
	admin.Post("/career-vacancy", c.CareerVacancyAdminController.Create)
	admin.Put("/career-vacancy/:id", c.CareerVacancyAdminController.Update)
	admin.Delete("/career-vacancy/:id", c.CareerVacancyAdminController.HardDelete)
}
