package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/database"
	"github.com/iyuz/devacademy-api/internal/functions/category"
	"github.com/iyuz/devacademy-api/internal/functions/course"
	"github.com/iyuz/devacademy-api/internal/functions/course/course_section"
	"github.com/iyuz/devacademy-api/internal/functions/user"
	"github.com/iyuz/devacademy-api/internal/routes"
)

type App struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Router *gin.Engine
}

func Init(cfg *config.Config) *App {
	db := database.InitPostgres(cfg)
	rdb := database.InitRedis(cfg)

	userService := user.NewUserService(user.NewUserRepository(db), cfg.JWT.Secret, cfg.JWT.Expiry)
	userController := user.NewUserController(userService)

	categoryService := category.NewCategoryService(category.NewCategoryRepository(db))
	categoryController := category.NewCategoryController(categoryService)

	courseService := course.NewCourseService(course.NewCourseRepository(db))
	courseController := course.NewCourseController(courseService)

	courseSectionService := coursesection.NewCourseSectionService(coursesection.NewCourseSectionRepository(db))
	courseSectionController := coursesection.NewCourseSectionController(courseSectionService, courseService)

	router := routes.SetupRouter(cfg, &routes.Controller{
		User:          userController,
		Category:      categoryController,
		Course:        courseController,
		CourseSection: courseSectionController,
	})

	return &App{
		DB:     db,
		Redis:  rdb,
		Router: router,
	}
}
