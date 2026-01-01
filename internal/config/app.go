package config

import (
	"aru-backend/internal/delivery/http"
	"aru-backend/internal/delivery/http/route"
	"aru-backend/internal/repository"
	"aru-backend/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validate    *validator.Validate
	Config      *viper.Viper
	MinioClient *minio.Client
	MinioConfig MinioConfig
}

func Bootstrap(config *BootstrapConfig) {

	// Repository layer
	// userRepository := repository.NewUserRepository()

	// Token service (JWT)
	// tokenService := service.NewTokenService()

	// Use case layer
	// userUseCase := usecase.NewUserUseCase(
	// 	config.DB,
	// 	config.Log,
	// 	config.Validate,
	// 	userRepository,
	// 	tokenService,
	// )

	// Controller layer
	// userController := http.NewUserController(
	// 	userUseCase,
	// 	config.Log,
	// )

	// Middleware layer
	// authMiddleware := middleware.NewAuthMiddleware(tokenService).Handle()

	// Route layer
	// routeConfig := route.RouteConfig{
	// 	App:            config.App,
	// 	UserController: userController,
	// 	AuthMiddleware: authMiddleware,
	// }

	// routeConfig.Setup()

	heroRepo := repository.NewHeroRepository(config.DB)

	siteContentUsecase := usecase.NewSiteContentUsecase(heroRepo, config.MinioConfig.PublicBaseURL)
	siteContentController := http.NewSiteContentController(siteContentUsecase)

	routeConfig := route.RouteConfig{
		App:                   config.App,
		SiteContentController: siteContentController,
	}

	routeConfig.Setup()
}
