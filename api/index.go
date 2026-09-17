package handler

import (
	"net/http"

	"ukm-hub/internal/config"
	"ukm-hub/internal/handler"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/routes"
	"ukm-hub/internal/service"

	"github.com/gin-gonic/gin"
)

var app *gin.Engine

func init() {
	db, err := config.InitDB()
	if err != nil {
		panic("Failed to connect DB: " + err.Error())
	}

	// 1. Repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	divRepo := repository.NewDivisionRepository(db)
	eventRepo := repository.NewEventRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	formRepo := repository.NewFormRepository(db)
	fieldRepo := repository.NewFormFieldRepository(db)
	memberRepo := repository.NewMemberRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// 2. Services
	accessService := service.NewAccessService(orgRepo)
	userService := service.NewUserService(userRepo, tokenRepo)
	orgService := service.NewOrganizationService(orgRepo, accessService)
	divService := service.NewDivisionService(divRepo, accessService)
	eventService := service.NewEventService(eventRepo, regRepo, accessService)
	formService := service.NewFormService(formRepo, fieldRepo, eventRepo, accessService)
	publicService := service.NewPublicService(eventRepo, formRepo, fieldRepo, orgRepo, regRepo)
	regService := service.NewRegistrationService(regRepo, eventRepo, formRepo, fieldRepo, accessService)
	memberService := service.NewMemberService(memberRepo, regRepo, eventRepo, formRepo, fieldRepo, divRepo, orgRepo, accessService)
	dashboardService := service.NewDashboardService(dashboardRepo, memberRepo, regRepo, formRepo, fieldRepo, divRepo, eventRepo, accessService)

	// 3. Handlers
	userHandler := handler.NewUserHandler(userService)
	orgHandler := handler.NewOrganizationHandler(orgService)
	divHandler := handler.NewDivisionHandler(divService)
	eventHandler := handler.NewEventHandler(eventService)
	formHandler := handler.NewFormHandler(formService)
	regHandler := handler.NewRegistrationHandler(regService)
	memberHandler := handler.NewMemberHandler(memberService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	publicHandler := handler.NewPublicHandler(publicService)

	gin.SetMode(gin.ReleaseMode)
	app = gin.New()
	app.Use(gin.Recovery())

	// 4. Passing 11 parameter ke SetupRouter
	routes.SetupRouter(
		app,
		userHandler,
		orgHandler,
		divHandler,
		eventHandler,
		formHandler,
		regHandler,
		memberHandler,
		dashboardHandler,
		publicHandler,
		tokenRepo,
	)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}