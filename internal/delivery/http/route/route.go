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

	HeroAdminController                 *http.HeroAdminController
	AwardAdminController                *http.AwardAdminController
	CareerVacancyAdminController        *http.CareerVacancyAdminController
	ClientAdminController               *http.ClientAdminController
	PartnerAdminController              *http.PartnerAdminController
	PromoSlideAdminController           *http.PromoSlideAdminController
	HistoryAdminController              *http.HistoryAdminController
	NewsCategoryAdminController         *http.NewsCategoryAdminController
	NewsArticleAdminController          *http.NewsArticleAdminController
	SpaceRoomAdminController            *http.SpaceRoomAdminController
	ServiceCertificationAdminController *http.ServiceCertificationAdminController
	ServiceMatrixAdminController        *http.ServiceMatrixAdminController
	UserAdminController                 *http.UserAdminController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Post("/api/auth/login", c.UserController.Login)
	c.App.Get("/api/me", c.UserController.Me)
	c.App.Put("/api/me", c.UserController.UpdateMe)
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

	admin.Get("/client", c.ClientAdminController.List)
	admin.Get("/client/:id", c.ClientAdminController.GetByID)
	admin.Post("/client", c.ClientAdminController.Create)
	admin.Put("/client/:id", c.ClientAdminController.Update)
	admin.Delete("/client/:id", c.ClientAdminController.HardDelete)

	admin.Get("/partner", c.PartnerAdminController.List)
	admin.Get("/partner/:id", c.PartnerAdminController.GetByID)
	admin.Post("/partner", c.PartnerAdminController.Create)
	admin.Put("/partner/:id", c.PartnerAdminController.Update)
	admin.Delete("/partner/:id", c.PartnerAdminController.HardDelete)

	admin.Get("/promo-slide", c.PromoSlideAdminController.List)
	admin.Get("/promo-slide/:id", c.PromoSlideAdminController.GetByID)
	admin.Post("/promo-slide", c.PromoSlideAdminController.Create)
	admin.Put("/promo-slide/:id", c.PromoSlideAdminController.Update)
	admin.Delete("/promo-slide/:id", c.PromoSlideAdminController.HardDelete)

	admin.Get("/history", c.HistoryAdminController.List)
	admin.Get("/history/:id", c.HistoryAdminController.GetByID)
	admin.Post("/history", c.HistoryAdminController.Create)
	admin.Put("/history/:id", c.HistoryAdminController.Update)
	admin.Delete("/history/:id", c.HistoryAdminController.HardDelete)

	admin.Get("/news-category", c.NewsCategoryAdminController.List)
	admin.Get("/news-category/:id", c.NewsCategoryAdminController.GetByID)
	admin.Post("/news-category", c.NewsCategoryAdminController.Create)
	admin.Put("/news-category/:id", c.NewsCategoryAdminController.Update)
	admin.Delete("/news-category/:id", c.NewsCategoryAdminController.HardDelete)

	admin.Get("/news-article", c.NewsArticleAdminController.List)
	admin.Get("/news-article/:id", c.NewsArticleAdminController.GetByID)
	admin.Post("/news-article", c.NewsArticleAdminController.Create)
	admin.Put("/news-article/:id", c.NewsArticleAdminController.Update)
	admin.Delete("/news-article/:id", c.NewsArticleAdminController.HardDelete)

	admin.Get("/space-room", c.SpaceRoomAdminController.List)
	admin.Get("/space-room/:id", c.SpaceRoomAdminController.GetByID)
	admin.Post("/space-room", c.SpaceRoomAdminController.Create)
	admin.Put("/space-room/:id", c.SpaceRoomAdminController.Update)
	admin.Delete("/space-room/:id", c.SpaceRoomAdminController.HardDelete)

	admin.Get("/service-certification/service/:service", c.ServiceCertificationAdminController.ListByService)
	admin.Get("/service-certification/:id", c.ServiceCertificationAdminController.GetByID)
	admin.Post("/service-certification", c.ServiceCertificationAdminController.Create)
	admin.Put("/service-certification/:id", c.ServiceCertificationAdminController.Update)
	admin.Delete("/service-certification/:id", c.ServiceCertificationAdminController.HardDelete)

	admin.Get("/service-matrix/service/:service", c.ServiceMatrixAdminController.ListByService)
	admin.Get("/service-matrix/:id", c.ServiceMatrixAdminController.GetByID)
	admin.Post("/service-matrix", c.ServiceMatrixAdminController.Create)
	admin.Put("/service-matrix/:id", c.ServiceMatrixAdminController.Update)
	admin.Delete("/service-matrix/:id", c.ServiceMatrixAdminController.HardDelete)

	admin.Get("/user", c.UserAdminController.List)
	admin.Get("/user/:id", c.UserAdminController.GetByID)
	admin.Post("/user", c.UserAdminController.Create)
	admin.Put("/user/:id", c.UserAdminController.Update)
	admin.Delete("/user/:id", c.UserAdminController.HardDelete)
}
