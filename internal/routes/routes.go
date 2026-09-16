package routes

import (
	"ukm-hub/internal/handler"
	"ukm-hub/internal/middleware"
	"ukm-hub/internal/repository"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler *handler.UserHandler, tokenRepo repository.TokenRepository) {
	api := r.Group("/api/v1")
	{
		// Public Routes (Bisa diakses tanpa login)
		auth := api.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", userHandler.Login)
			auth.POST("/logout", middleware.AuthMiddleware(tokenRepo), userHandler.Logout)
		}

		// Protected Routes (Wajib membawa Token JWT di Header)
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware(tokenRepo))
		{
			users.GET("/me", userHandler.GetProfile)
			users.PUT("/me", userHandler.UpdateProfile)
		}
	}
}
