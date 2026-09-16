package main

import (
	"log"
	"path/filepath"
	"runtime"

	"ukm-hub/internal/config"
	"ukm-hub/internal/entity"
	"ukm-hub/internal/handler"
	"ukm-hub/internal/repository"
	"ukm-hub/internal/routes"
	"ukm-hub/internal/service"

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
	db.AutoMigrate(&entity.User{}, &entity.RevokedToken{})

	// Dependency Injection Wiring
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	tokenRepo.DeleteExpired()

	userService := service.NewUserService(userRepo, tokenRepo)
	userHandler := handler.NewUserHandler(userService)

	if config.Get("ENV", "development") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Setup Routes
	routes.SetupRouter(r, userHandler, tokenRepo)

	r.Run(":" + config.Get("PORT", "8080"))
}

func loadDotEnv() {
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(file))
	_ = godotenv.Load(filepath.Join(projectRoot, ".env"))
	_ = godotenv.Load(".env")
}
