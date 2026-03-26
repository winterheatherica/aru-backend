package config

import (
	"aru-backend/internal/delivery/http"
	"aru-backend/internal/delivery/http/middleware"
	"aru-backend/internal/delivery/http/route"
	"aru-backend/internal/repository"
	"aru-backend/internal/usecase"
	"aru-backend/internal/usecase/admin"

	"time"

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
	newsRepo := repository.NewNewsArticleRepository(config.DB)

	homeUsecase := usecase.NewHomeUsecase(
		heroRepo,
		promoRepo,
		partnerRepo,
		clientRepo,
		newsRepo,
		config.MinioConfig.PublicBaseURL,
	)

	homeController := http.NewHomeController(homeUsecase)

	// --- About Module ---
	historyRepo := repository.NewHistoryRepository(config.DB)
	awardRepo := repository.NewAwardRepository(config.DB)

	aboutUsecase := usecase.NewAboutUsecase(
		historyRepo,
		partnerRepo,
		awardRepo,
		config.MinioConfig.PublicBaseURL,
	)

	aboutController := http.NewAboutController(aboutUsecase)

	// --- Service Module ---
	serviceGalleryRepo := repository.NewServiceGalleryRepository(config.DB)
	servicePricingRepo := repository.NewServicePricingRepository(config.DB)
	serviceMatrixRepo := repository.NewServiceMatrixRepository(config.DB)
	serviceCertificationRepo := repository.NewServiceCertificationRepository(config.DB)

	serviceUsecase := usecase.NewServiceUsecase(
		serviceGalleryRepo,
		servicePricingRepo,
		serviceMatrixRepo,
		serviceCertificationRepo,
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

	// --- User/Auth Module ---
	userRepo := repository.NewUserRepository(config.DB)
	sessionRepo := repository.NewSessionRepository(config.DB)
	jwtSecret := config.Config.GetString("jwt.secret")
	jwtExpiry := time.Duration(config.Config.GetInt("jwt.expiryHours")) * time.Hour
	if jwtExpiry <= 0 {
		jwtExpiry = 24 * time.Hour
	}
	userUsecase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepo, sessionRepo, jwtSecret, jwtExpiry)
	userController := http.NewUserController(userUsecase, config.Log)

	// --- Admin Hero Module ---
	heroAdminUsecase := admin.NewHeroUsecase(
		heroRepo,
		config.MinioClient,
		config.MinioConfig.Bucket,
		config.MinioConfig.PublicBaseURL,
	)
	heroAdminController := http.NewHeroAdminController(heroAdminUsecase)

	// --- Admin Award Module ---
	awardAdminUsecase := admin.NewAwardUsecase(
		awardRepo,
		config.MinioClient,
		config.MinioConfig.Bucket,
		config.MinioConfig.PublicBaseURL,
	)
	awardAdminController := http.NewAwardAdminController(awardAdminUsecase)

	// --- Admin Career Vacancy Module ---
	careerVacancyAdminUsecase := admin.NewCareerVacancyUsecase(careerVacancyRepo)
	careerVacancyAdminController := http.NewCareerVacancyAdminController(careerVacancyAdminUsecase)

	// --- Admin Client Module ---
	clientAdminUsecase := admin.NewClientUsecase(
		clientRepo,
		config.MinioClient,
		config.MinioConfig.Bucket,
		config.MinioConfig.PublicBaseURL,
	)
	clientAdminController := http.NewClientAdminController(clientAdminUsecase)

	authMiddleware := middleware.NewAuthMiddleware()

	// --- Setup Routes ---
	routeConfig := route.RouteConfig{
		App: config.App,

		AuthMiddleware: authMiddleware.Handle(),

		UserController: userController,

		HomeController:        homeController,
		AboutController:       aboutController,
		ServiceController:     serviceController,
		ReservationController: reservationController,
		InformationController: informationController,
		CareerController:      careerController,

		RoomController: roomController,

		ArticleController:  articleController,
		CategoryController: categoryController,

		HeroAdminController:          heroAdminController,
		AwardAdminController:         awardAdminController,
		CareerVacancyAdminController: careerVacancyAdminController,
		ClientAdminController:        clientAdminController,
	}

	routeConfig.Setup()

}
