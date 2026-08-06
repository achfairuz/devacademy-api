package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/iyuz/devacademy-api/internal/config"
	"github.com/iyuz/devacademy-api/internal/models"
)

func InitPostgres(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if viper.GetBool("DB_AUTOMIGRATE") {
		if err := db.AutoMigrate(
			&models.User{},
			&models.Category{},
			&models.Course{},
			&models.CourseSection{},
			&models.Lesson{},
			&models.LessonFile{},
			&models.Enrollment{},
			&models.LessonProgress{},
			&models.Quiz{},
			&models.QuizQuestion{},
			&models.QuizOption{},
			&models.QuizAttempt{},
			&models.Assignment{},
			&models.AssignmentSubmission{},
			&models.SubscriptionPlan{},
			&models.UserSubscription{},
			&models.Payment{},
			&models.Certificate{},
		); err != nil {
			log.Fatalf("failed to run migration: %v", err)
		}
	}

	return db
}

func InitRedis(cfg *config.Config) *redis.Client {
	if !cfg.Redis.Enabled {
		log.Println("redis is disabled, skipping connection")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	return client
}
