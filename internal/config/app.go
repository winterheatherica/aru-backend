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

	// --- Service Module ---
	serviceGalleryRepo := repository.NewServiceGalleryRepository(config.DB)
	servicePricingRepo := repository.NewServicePricingRepository(config.DB)
	serviceMatrixRepo := repository.NewServiceMatrixRepository(config.DB)

	serviceUsecase := usecase.NewServiceUsecase(
		serviceGalleryRepo,
		servicePricingRepo,
		serviceMatrixRepo,
		config.MinioConfig.PublicBaseURL,
	)

	serviceController := http.NewServiceController(serviceUsecase)

	// --- Reservation Module ---
	spaceRoomRepo := repository.NewSpaceRoomRepository(config.DB)

	reservationUsecase := usecase.NewReservationUsecase(
		spaceRoomRepo,
		config.MinioConfig.PublicBaseURL,
	)

	reservationController := http.NewReservationController(reservationUsecase)

	// --- Information Module ---
	newsArticleRepo := repository.NewNewsArticleRepository(config.DB)

	informationUsecase := usecase.NewInformationUsecase(
		newsArticleRepo,
		config.MinioConfig.PublicBaseURL,
	)

	informationController := http.NewInformationController(informationUsecase)

	// --- Career Module ---
	careerVacancyRepo := repository.NewCareerVacancyRepository(config.DB)

	careerUsecase := usecase.NewCareerUsecase(
		careerVacancyRepo,
	)

	careerController := http.NewCareerController(careerUsecase)

	// --- Room Module ---
	roomDetailRepo := repository.NewSpaceRoomRepository(config.DB)

	roomUsecase := usecase.NewRoomUsecase(
		roomDetailRepo,
		config.MinioConfig.PublicBaseURL,
	)

	roomController := http.NewRoomController(roomUsecase)

	// --- Article Module ---
	articleRepo := repository.NewNewsArticleRepository(config.DB)

	articleUsecase := usecase.NewArticleUsecase(
		articleRepo,
		config.MinioConfig.PublicBaseURL,
	)

	articleController := http.NewArticleController(articleUsecase)

	// --- Category Module ---
	categoryRepo := repository.NewNewsCategoryRepository(config.DB)

	categoryUsecase := usecase.NewCategoryUsecase(
		categoryRepo,
	)

	categoryController := http.NewCategoryController(categoryUsecase)

	// --- Setup Routes ---
	routeConfig := route.RouteConfig{
		App: config.App,

		HomeController:        homeController,
		AboutController:       aboutController,
		ServiceController:     serviceController,
		ReservationController: reservationController,
		InformationController: informationController,
		CareerController:      careerController,

		RoomController: roomController,

		ArticleController:  articleController,
		CategoryController: categoryController,
	}

	routeConfig.Setup()

}
