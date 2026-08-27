package main

import (
	"log"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/iyuz/devacademy-api/docs"
	"github.com/iyuz/devacademy-api/internal/bootstrap"
	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/database"
)

// @title			DevAcademy API
// @version			1.0
// @description		API documentation for DevAcademy backend service.
// @host			localhost:8080
// @BasePath		/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in								header
// @name							Authorization
func main() {
	cfg := config.Load()

	app := bootstrap.Init(cfg)

	if err := database.SeedUsers(app.DB); err != nil {
		log.Fatalf("failed to seed users: %v", err)
	}

	if err := database.SeedLevels(app.DB); err != nil {
		log.Fatalf("failed to seed levels: %v", err)
	}

	if err := database.SeedCourseContent(app.DB); err != nil {
		log.Fatalf("failed to seed course content: %v", err)
	}

	if err := database.SeedEnrollments(app.DB); err != nil {
		log.Fatalf("failed to seed enrollments: %v", err)
	}

	app.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Printf("server running on port %s", cfg.Server.Port)

	println("Swagger documentation: http://localhost:" + cfg.Server.Port + "/swagger/index.html")
	if err := app.Router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
