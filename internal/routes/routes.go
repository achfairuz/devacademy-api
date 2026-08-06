package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/handlers"
	"github.com/iyuz/devacademy-api/internal/middleware"
	"github.com/iyuz/devacademy-api/internal/models"
)

type Handler struct {
	User *handlers.UserHandler
}

func SetupRouter(cfg *config.Config, h *Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", h.User.Register)
			auth.POST("/login", h.User.Login)
		}

		users := api.Group("/users", middleware.Auth(cfg.JWT.Secret))
		{
			users.GET("", middleware.RequireRole(models.RoleAdmin), h.User.GetAll)
			users.GET("/:id", h.User.GetByID)
			users.PUT("/:id", h.User.Update)
			users.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), h.User.Delete)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
