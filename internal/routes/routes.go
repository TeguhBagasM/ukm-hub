package routes

import (
	"handler/internal/handler"
	"handler/internal/middleware"
	"handler/internal/repository"
	"handler/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, userHandler *handler.UserHandler, orgHandler *handler.OrganizationHandler, divHandler *handler.DivisionHandler, eventHandler *handler.EventHandler, formHandler *handler.FormHandler, regHandler *handler.RegistrationHandler, memberHandler *handler.MemberHandler, dashboardHandler *handler.DashboardHandler, publicHandler *handler.PublicHandler, tokenRepo repository.TokenRepository) {
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

		// Organization Management (SUPER_ADMIN semua, ORG_ADMIN terbatas ke organisasinya)
		orgs := api.Group("/organizations")
		orgs.Use(middleware.AuthMiddleware(tokenRepo))
		{
			orgs.GET("", middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), orgHandler.List)
			orgs.POST("", middleware.RequireRole(utils.RoleSuperAdmin), orgHandler.Create)
			orgs.GET("/:id", middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), orgHandler.Get)
			orgs.PUT("/:id", middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), orgHandler.Update)
			orgs.DELETE("/:id", middleware.RequireRole(utils.RoleSuperAdmin), orgHandler.Delete)
		}

		// Division Management (scoped ke organisasi)
		api.GET("/organizations/:id/divisions", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), divHandler.ListByOrganization)
		api.POST("/organizations/:id/divisions", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), divHandler.Create)
		api.GET("/divisions/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), divHandler.Get)
		api.PUT("/divisions/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), divHandler.Update)
		api.DELETE("/divisions/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), divHandler.Delete)

		// Event Management (scoped ke organisasi)
		api.GET("/organizations/:id/events", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.ListByOrganization)
		api.POST("/organizations/:id/events", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Create)
		api.GET("/events/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Get)
		api.PUT("/events/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Update)
		api.DELETE("/events/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Delete)
		api.POST("/events/:id/publish", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Publish)
		api.POST("/events/:id/close", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Close)
		api.POST("/events/:id/archive", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), eventHandler.Archive)

		// Dynamic Form Builder (scoped ke organisasi)
		api.GET("/events/:id/form", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.GetByEvent)
		api.POST("/events/:id/form", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.CreateForm)
		api.PUT("/forms/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.UpdateForm)
		api.POST("/forms/:id/publish", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.Publish)
		api.POST("/forms/:id/unpublish", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.Unpublish)
		api.POST("/forms/:id/fields", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.AddField)
		api.PUT("/fields/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.UpdateField)
		api.DELETE("/fields/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.DeleteField)
		api.POST("/forms/:id/fields/:fieldId/duplicate", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.DuplicateField)
		api.POST("/forms/:id/fields/reorder", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), formHandler.ReorderFields)

		// Registration Management (scoped ke organisasi)
		api.GET("/events/:id/registrations", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), regHandler.ListByEvent)
		api.GET("/registrations/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), regHandler.Get)
		api.POST("/registrations/:id/accept", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), regHandler.Accept)
		api.POST("/registrations/:id/reject", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), regHandler.Reject)

		// Member Management (scoped ke organisasi)
		api.GET("/organizations/:id/members", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), memberHandler.List)
		api.POST("/registrations/:id/convert-member", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), memberHandler.Convert)
		api.GET("/members/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), memberHandler.Get)
		api.PUT("/members/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), memberHandler.Update)
		api.DELETE("/members/:id", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), memberHandler.Delete)

		// Dashboard & CSV Export (scoped ke organisasi)
		api.GET("/organizations/:id/dashboard", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), dashboardHandler.Get)
		api.GET("/organizations/:id/registrations/export", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), dashboardHandler.ExportRegistrations)
		api.GET("/organizations/:id/members/export", middleware.AuthMiddleware(tokenRepo), middleware.RequireRole(utils.RoleSuperAdmin, utils.RoleOrgAdmin), dashboardHandler.ExportMembers)

		// Public Registration (tanpa login)
		public := api.Group("/public")
		{
			public.GET("/events/:slug", publicHandler.GetEvent)
			public.GET("/events/:slug/form", publicHandler.GetForm)
			public.POST("/events/:slug/registrations", publicHandler.Register)
		}
	}
}
