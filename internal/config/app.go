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

	// --- Home Module ---
	heroRepo := repository.NewHeroRepository(config.DB)
	promoRepo := repository.NewPromoRepository(config.DB)
	partnerRepo := repository.NewPartnerRepository(config.DB)
	clientRepo := repository.NewClientRepository(config.DB)

	homeUsecase := usecase.NewHomeUsecase(
		heroRepo,
		promoRepo,
		partnerRepo,
		clientRepo,
		config.MinioConfig.PublicBaseURL,
	)

	homeController := http.NewHomeController(homeUsecase)

	// --- About Module ---
	historyRepo := repository.NewHistoryRepository(config.DB)
	awardRepo := repository.NewAwardRepository(config.DB)

	aboutUsecase := usecase.NewAboutUsecase(
		historyRepo,
		awardRepo,
		config.MinioConfig.PublicBaseURL,
	)

	aboutController := http.NewAboutController(aboutUsecase)

	// --- Career Module ---
	careerVacancyRepo := repository.NewCareerVacancyRepository(config.DB)

	careerUsecase := usecase.NewCareerUsecase(
		careerVacancyRepo,
	)

	careerController := http.NewCareerController(careerUsecase)

	// --- Setup Routes ---
	routeConfig := route.RouteConfig{
		App:              config.App,
		HomeController:   homeController,
		AboutController:  aboutController,
		CareerController: careerController,
	}

	routeConfig.Setup()

}
