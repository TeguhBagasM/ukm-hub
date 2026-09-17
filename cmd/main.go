package main

import (
	"log"
	"path/filepath"
	"runtime"

	"handler/internal/config"
	"handler/internal/entity"
	"handler/internal/handler"
	"handler/internal/middleware"
	"handler/internal/repository"
	"handler/internal/routes"
	"handler/internal/seed"
	"handler/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadDotEnv()

	// Init DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// Auto Migrate Models
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.RevokedToken{},
		&entity.Organization{},
		&entity.OrganizationAdmin{},
		&entity.Division{},
		&entity.Event{},
		&entity.Form{},
		&entity.FormField{},
		&entity.Registration{},
		&entity.RegistrationAnswer{},
		&entity.Member{},
	); err != nil {
		log.Fatalf("Failed to migrate DB: %v", err)
	}

	// Seed Demo Data (Idempotent)
	if err := seed.Run(db); err != nil {
		log.Fatalf("Failed to seed DB: %v", err)
	}

	// Dependency Injection Wiring
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	tokenRepo.DeleteExpired()

	userService := service.NewUserService(userRepo, tokenRepo)
	userHandler := handler.NewUserHandler(userService)

	orgRepo := repository.NewOrganizationRepository(db)
	accessService := service.NewAccessService(orgRepo)
	orgService := service.NewOrganizationService(orgRepo, accessService)
	orgHandler := handler.NewOrganizationHandler(orgService)

	divRepo := repository.NewDivisionRepository(db)
	divService := service.NewDivisionService(divRepo, accessService)
	divHandler := handler.NewDivisionHandler(divService)

	eventRepo := repository.NewEventRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	eventService := service.NewEventService(eventRepo, regRepo, accessService)
	eventHandler := handler.NewEventHandler(eventService)

	formRepo := repository.NewFormRepository(db)
	fieldRepo := repository.NewFormFieldRepository(db)
	formService := service.NewFormService(formRepo, fieldRepo, eventRepo, accessService)
	formHandler := handler.NewFormHandler(formService)

	publicService := service.NewPublicService(eventRepo, formRepo, fieldRepo, orgRepo, regRepo)
	publicHandler := handler.NewPublicHandler(publicService)

	regService := service.NewRegistrationService(regRepo, eventRepo, formRepo, fieldRepo, accessService)
	regHandler := handler.NewRegistrationHandler(regService)

	memberRepo := repository.NewMemberRepository(db)
	memberService := service.NewMemberService(memberRepo, regRepo, eventRepo, formRepo, fieldRepo, divRepo, orgRepo, accessService)
	memberHandler := handler.NewMemberHandler(memberService)

	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepo, memberRepo, regRepo, formRepo, fieldRepo, divRepo, eventRepo, accessService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)

	if config.Get("ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Setup Routes
	routes.SetupRouter(r, userHandler, orgHandler, divHandler, eventHandler, formHandler, regHandler, memberHandler, dashboardHandler, publicHandler, tokenRepo)

	r.Run(":" + config.Get("PORT", "8080"))
}

func loadDotEnv() {
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(file))
	_ = godotenv.Load(filepath.Join(projectRoot, ".env"))
	_ = godotenv.Load(".env")
}
