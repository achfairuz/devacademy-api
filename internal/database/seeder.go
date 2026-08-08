package database

import (
	"errors"
	"log"
	"time"

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

func SeedLevels(db *gorm.DB) error {
	levels := []models.Level{
		{Name: "Beginner", Slug: "beginner"},
		{Name: "Intermediate", Slug: "intermediate"},
		{Name: "Advanced", Slug: "advanced"},
	}

	for _, level := range levels {
		var existing models.Level
		err := db.Where("slug = ?", level.Slug).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := db.Create(&level).Error; err != nil {
			return err
		}
	}

	log.Println("levels seeded successfully")
	return nil
}

func SeedCourseContent(db *gorm.DB) error {
	var mentor models.User
	if err := db.Where("email = ?", "mentor@devacademy.com").First(&mentor).Error; err != nil {
		return err
	}

	category := models.Category{Name: "Backend Development", Slug: "backend-development", Icon: "server"}
	if err := db.Where("slug = ?", category.Slug).FirstOrCreate(&category).Error; err != nil {
		return err
	}

	var level models.Level
	if err := db.Where("slug = ?", "beginner").First(&level).Error; err != nil {
		return err
	}

	var existing models.Course
	err := db.Where("slug = ?", "belajar-golang-dari-dasar").First(&existing).Error
	if err == nil {
		log.Println("course content already seeded, skipping")
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	course := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  category.ID,
		LevelID:     level.ID,
		Title:       "Belajar Golang dari Dasar",
		Slug:        "belajar-golang-dari-dasar",
		Description: "Belajar bahasa pemrograman Go mulai dari nol hingga siap membuat aplikasi.",
		Price:       150000,
		Duration:    12,
		Status:      "published",
	}

	sectionDasar := models.CourseSection{Title: "Pengenalan Golang", OrderNumber: 1}
	sectionLanjutan := models.CourseSection{Title: "Struktur Data & Concurrency", OrderNumber: 2}

	lesson1 := models.Lesson{
		Title:       "Apa itu Golang",
		Description: "Mengenal sejarah dan keunggulan bahasa Go.",
		VideoURL:    "https://www.youtube.com/watch?v=demo1",
		OrderNumber: 1,
		IsPreview:   true,
	}
	lesson2 := models.Lesson{
		Title:       "Instalasi & Setup",
		Description: "Cara instalasi Go dan setup environment.",
		VideoURL:    "https://www.youtube.com/watch?v=demo2",
		OrderNumber: 2,
	}
	lesson3 := models.Lesson{
		Title:       "Variabel & Tipe Data",
		Description: "Dasar variabel dan tipe data di Golang.",
		VideoURL:    "https://www.youtube.com/watch?v=demo3",
		OrderNumber: 3,
	}
	lesson4 := models.Lesson{
		Title:       "Slice, Map & Struct",
		Description: "Koleksi data dan struct di Golang.",
		VideoURL:    "https://www.youtube.com/watch?v=demo4",
		OrderNumber: 4,
	}
	lesson5 := models.Lesson{
		Title:       "Goroutine & Channel",
		Description: "Concurrency di Golang menggunakan goroutine dan channel.",
		VideoURL:    "https://www.youtube.com/watch?v=demo5",
		OrderNumber: 5,
	}

	quiz := &models.Quiz{
		Title:        "Kuis Dasar Golang",
		PassingScore: 60,
		Questions: []models.QuizQuestion{
			{
				Question:     "Siapa penemu bahasa pemrograman Go?",
				QuestionType: "multiple_choice",
				Options: []models.QuizOption{
					{OptionText: "Google", IsCorrect: true},
					{OptionText: "Microsoft", IsCorrect: false},
					{OptionText: "Facebook", IsCorrect: false},
					{OptionText: "Apple", IsCorrect: false},
				},
			},
			{
				Question:     "Keyword apa yang digunakan untuk mendeklarasikan variabel yang nilainya tidak dapat diubah?",
				QuestionType: "multiple_choice",
				Options: []models.QuizOption{
					{OptionText: "var", IsCorrect: false},
					{OptionText: "let", IsCorrect: false},
					{OptionText: "const", IsCorrect: true},
					{OptionText: "static", IsCorrect: false},
				},
			},
			{
				Question:     "Apakah Golang memiliki garbage collector?",
				QuestionType: "true_false",
				Options: []models.QuizOption{
					{OptionText: "True", IsCorrect: true},
					{OptionText: "False", IsCorrect: false},
				},
			},
		},
	}

	dueDate := time.Now().Add(7 * 24 * time.Hour)
	assignment := &models.Assignment{
		Title:       "Latihan Mandiri: Hello World",
		Description: "Buat program sederhana yang mencetak 'Hello, World!' dan jelaskan cara kerjanya.",
		DueDate:     &dueDate,
	}

	lesson1.Quiz = quiz
	lesson2.Assignment = assignment
	sectionDasar.Lessons = []models.Lesson{lesson1, lesson2, lesson3}
	sectionLanjutan.Lessons = []models.Lesson{lesson4, lesson5}
	course.Sections = []models.CourseSection{sectionDasar, sectionLanjutan}

	if err := db.Create(course).Error; err != nil {
		return err
	}

	log.Println("course content seeded successfully")
	return nil
}
