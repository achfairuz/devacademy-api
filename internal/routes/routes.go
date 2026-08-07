package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/controllers"
	"github.com/iyuz/devacademy-api/internal/middleware"
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

		admin := api.Group("/admin", middleware.Auth(cfg.JWT.Secret), middleware.AdminOnly())
		{
			admin.GET("/users", ctr.User.GetAll)
			admin.GET("/users/:id", ctr.User.GetByID)
			admin.PUT("/users/:id", ctr.User.Update)
			admin.DELETE("/users/:id", ctr.User.Delete)
		}

		mentor := api.Group("/mentor", middleware.Auth(cfg.JWT.Secret), middleware.MentorOnly())
		{
			mentor.GET("/categories", ctr.Category.GetAll)
			mentor.GET("/categories/:id", ctr.Category.GetByID)
			mentor.POST("/categories", ctr.Category.Create)
			mentor.DELETE("/categories/:id", ctr.Category.Delete)
		}

		student := api.Group("/student", middleware.Auth(cfg.JWT.Secret), middleware.StudentOnly())
		{
			student.GET("/categories", ctr.Category.GetAll)
			student.GET("/categories/:id", ctr.Category.GetByID)
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
