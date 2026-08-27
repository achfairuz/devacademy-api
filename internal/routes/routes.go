package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/functions/category"
	"github.com/iyuz/devacademy-api/internal/functions/course"
	coursesection "github.com/iyuz/devacademy-api/internal/functions/course/course_section"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons/assignment"
	lessonfile "github.com/iyuz/devacademy-api/internal/functions/course/lessons/lesson_file"
	"github.com/iyuz/devacademy-api/internal/functions/course/lessons/quiz"
	"github.com/iyuz/devacademy-api/internal/functions/level"
	"github.com/iyuz/devacademy-api/internal/functions/user"
	"github.com/iyuz/devacademy-api/internal/middleware"
	"github.com/iyuz/devacademy-api/internal/models"
)

type Controller struct {
	User          *user.UserController
	Category      *category.CategoryController
	Course        *course.CourseController
	CourseSection *coursesection.CourseSectionController
	Lesson        *lessons.LessonController
	LessonFile    *lessonfile.LessonFileController
	Quiz          *quiz.QuizController
	Assignment    *assignment.AssignmentController
	Level         *level.LevelController
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

		levels := api.Group("/levels")
		{
			levels.GET("", ctr.Level.GetAll)
			levels.GET("/:id", ctr.Level.GetByID)
			levels.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleAdmin), ctr.Level.Create)
			levels.PUT("/:id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleAdmin), ctr.Level.Update)
			levels.DELETE("/:id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleAdmin), ctr.Level.Delete)
		}

		courses := api.Group("/courses")
		{
			courses.GET("", ctr.Course.GetAll)
			courses.GET("/cards", middleware.OptionalAuth(cfg.JWT.Secret), ctr.Course.GetCards)
			courses.GET("/:id", ctr.Course.GetByID)
			courses.GET("/slug/:slug", ctr.Course.GetBySlug)
			courses.GET("/mentor/:mentor_id", ctr.Course.GetByMentor)
			courses.GET("/category/:category_id", ctr.Course.GetByCategory)
			courses.GET("/level/:level_id", ctr.Course.GetByLevel)

			courses.GET("/slug/:slug/detail", middleware.OptionalAuth(cfg.JWT.Secret), ctr.Course.GetDetailBySlug)

			courses.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Course.Create)
			courses.PUT("/:id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Course.Update)
			courses.PATCH("/slug/:slug/status", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Course.UpdateStatus)
			courses.DELETE("/:id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Course.Delete)

			courseSections := courses.Group("/:id/sections")
			{
				courseSections.GET("", ctr.CourseSection.GetByCourse)
				courseSections.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.CourseSection.Create)
				courseSections.GET("/:section_id", ctr.CourseSection.GetByID)
				courseSections.PUT("/:section_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.CourseSection.Update)
				courseSections.DELETE("/:section_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.CourseSection.Delete)

				courseLessons := courseSections.Group("/:section_id/lessons")
				{
					courseLessons.GET("", ctr.Lesson.Index)
					courseLessons.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Lesson.Store)
					courseLessons.GET("/:lesson_id", ctr.Lesson.Show)
					courseLessons.PUT("/:lesson_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Lesson.Update)
					courseLessons.DELETE("/:lesson_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Lesson.Delete)

					courseFiles := courseLessons.Group("/:lesson_id/files")
					{
						courseFiles.GET("", ctr.LessonFile.Index)
						courseFiles.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.LessonFile.Store)
						courseFiles.GET("/:file_id", ctr.LessonFile.Show)
						courseFiles.DELETE("/:file_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.LessonFile.Delete)
					}

					courseQuizzes := courseLessons.Group("/:lesson_id/quizzes")
					{
						courseQuizzes.GET("", ctr.Quiz.Index)
						courseQuizzes.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Quiz.Store)
						courseQuizzes.GET("/:quiz_id", ctr.Quiz.Show)
						courseQuizzes.PUT("/:quiz_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Quiz.Update)
						courseQuizzes.DELETE("/:quiz_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Quiz.Delete)
					}

					courseAssignments := courseLessons.Group("/:lesson_id/assignments")
					{
						courseAssignments.GET("", ctr.Assignment.Index)
						courseAssignments.POST("", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Assignment.Store)
						courseAssignments.GET("/:assignment_id", ctr.Assignment.Show)
						courseAssignments.PUT("/:assignment_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Assignment.Update)
						courseAssignments.DELETE("/:assignment_id", middleware.Auth(cfg.JWT.Secret), middleware.RequireRole(models.RoleMentor, models.RoleAdmin), ctr.Assignment.Delete)
					}
				}
			}
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
