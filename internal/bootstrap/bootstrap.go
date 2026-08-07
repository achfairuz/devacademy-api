package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/controllers"
	"github.com/iyuz/devacademy-api/internal/database"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"github.com/iyuz/devacademy-api/internal/repositories/impl"
	"github.com/iyuz/devacademy-api/internal/routes"
	"github.com/iyuz/devacademy-api/internal/services"
)

type App struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Router *gin.Engine
}

func Init(cfg *config.Config) *App {
	db := database.InitPostgres(cfg)
	rdb := database.InitRedis(cfg)

	userRepo := impl.NewUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWT.Secret, cfg.JWT.Expiry)
	userController := controllers.NewUserController(userService)

	categoryService := services.NewCategoryService(repositories.NewCategoryRepository(db))
	categoryController := controllers.NewCategoryController(categoryService)

	courseService := services.NewCourseService(impl.NewCourseRepository(db))
	courseController := controllers.NewCourseController(courseService)

	router := routes.SetupRouter(cfg, &routes.Controller{
		User:     userController,
		Category: categoryController,
		Course:   courseController,
	})

	return &App{
		DB:     db,
		Redis:  rdb,
		Router: router,
	}
}
