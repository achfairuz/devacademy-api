package database

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/utils"
)

func SeedUsers(db *gorm.DB) error {
	users := []models.User{
		{
			FullName:      "Admin DevAcademy",
			Username:      "admin",
			Email:         "admin@devacademy.com",
			Password:      "admin123",
			Role:          models.RoleAdmin,
			Status:        true,
			EmailVerified: true,
		},
		{
			FullName:      "Mentor DevAcademy",
			Username:      "mentor",
			Email:         "mentor@devacademy.com",
			Password:      "mentor123",
			Role:          models.RoleMentor,
			Status:        true,
			EmailVerified: true,
		},
		{
			FullName:      "Student DevAcademy",
			Username:      "student",
			Email:         "student@devacademy.com",
			Password:      "student123",
			Role:          models.RoleStudent,
			Status:        true,
			EmailVerified: true,
		},
	}

	for _, user := range users {
		var existing models.User
		err := db.Where("email = ?", user.Email).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		hash, err := utils.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hash

		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	log.Println("users seeded successfully")
	return nil
}
