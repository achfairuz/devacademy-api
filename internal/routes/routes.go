package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/controllers"
	"github.com/iyuz/devacademy-api/internal/middleware"
	"github.com/iyuz/devacademy-api/internal/models"
)

type Controller struct {
	User     *controllers.UserController
	Category *controllers.CategoryController
}

func SetupRouter(cfg *config.Config, ctr *Controller) *gin.Engine {
	router := gin.New()
	router.MaxMultipartMemory = 8 << 20
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORSMiddleware())

	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", ctr.User.Register)
			auth.POST("/login", ctr.User.Login)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", ctr.Category.GetAll)
			categories.GET("/:id", ctr.Category.GetByID)
			categories.POST("", ctr.Category.Create, middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleAdmin))
			categories.DELETE("/:id", ctr.Category.Delete, middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleAdmin))
		}

		users := api.Group("/users", middleware.Auth(cfg.JWT.Secret))
		{
			users.GET("", middleware.RequireRole(models.RoleAdmin), ctr.User.GetAll)
			users.GET("/:id", ctr.User.GetByID)
			users.PUT("/:id", ctr.User.Update)
			users.DELETE("/:id", middleware.RequireRole(models.RoleAdmin), ctr.User.Delete)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
