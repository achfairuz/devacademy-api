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

	categoryBackend := models.Category{Name: "Backend Development", Slug: "backend-development", Icon: "server"}
	if err := db.Where("slug = ?", categoryBackend.Slug).FirstOrCreate(&categoryBackend).Error; err != nil {
		return err
	}

	categoryFrontend := models.Category{Name: "Frontend Development", Slug: "frontend-development", Icon: "layout"}
	if err := db.Where("slug = ?", categoryFrontend.Slug).FirstOrCreate(&categoryFrontend).Error; err != nil {
		return err
	}

	categoryData := models.Category{Name: "Data Science", Slug: "data-science", Icon: "bar-chart"}
	if err := db.Where("slug = ?", categoryData.Slug).FirstOrCreate(&categoryData).Error; err != nil {
		return err
	}

	categoryMobile := models.Category{Name: "Mobile Development", Slug: "mobile-development", Icon: "smartphone"}
	if err := db.Where("slug = ?", categoryMobile.Slug).FirstOrCreate(&categoryMobile).Error; err != nil {
		return err
	}

	levelBeginner := models.Level{}
	if err := db.Where("slug = ?", "beginner").First(&levelBeginner).Error; err != nil {
		return err
	}

	levelIntermediate := models.Level{}
	if err := db.Where("slug = ?", "intermediate").First(&levelIntermediate).Error; err != nil {
		return err
	}

	levelAdvanced := models.Level{}
	if err := db.Where("slug = ?", "advanced").First(&levelAdvanced).Error; err != nil {
		return err
	}

	course1 := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  categoryBackend.ID,
		LevelID:     levelBeginner.ID,
		Title:       "Belajar Golang dari Dasar",
		Slug:        "belajar-golang-dari-dasar",
		Description: "Belajar bahasa pemrograman Go mulai dari nol hingga siap membuat aplikasi.",
		Price:       150000,
		Duration:    12,
		Status:      "published",
		Sections: []models.CourseSection{
			{
				Title:       "Pengenalan Golang",
				OrderNumber: 1,
				Lessons: []models.Lesson{
					{
						Title:       "Apa itu Golang",
						Description: "Mengenal sejarah dan keunggulan bahasa Go.",
						VideoURL:    "https://www.youtube.com/watch?v=demo1",
						Duration:    600,
						OrderNumber: 1,
						IsPreview:   true,
						Quiz: &models.Quiz{
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
									Question:     "Keyword untuk mendeklarasikan variabel yang tidak dapat diubah?",
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
						},
					},
					{
						Title:       "Instalasi & Setup",
						Description: "Cara instalasi Go dan setup environment.",
						VideoURL:    "https://www.youtube.com/watch?v=demo2",
						Duration:    480,
						OrderNumber: 2,
						Assignment: &models.Assignment{
							Title:       "Latihan Mandiri: Hello World",
							Description: "Buat program sederhana yang mencetak 'Hello, World!' dan jelaskan cara kerjanya.",
							DueDate:     timePtr(time.Now().Add(7 * 24 * time.Hour)),
						},
					},
					{
						Title:       "Variabel & Tipe Data",
						Description: "Dasar variabel dan tipe data di Golang.",
						VideoURL:    "https://www.youtube.com/watch?v=demo3",
						Duration:    720,
						OrderNumber: 3,
					},
				},
			},
			{
				Title:       "Struktur Data & Concurrency",
				OrderNumber: 2,
				Lessons: []models.Lesson{
					{
						Title:       "Slice, Map & Struct",
						Description: "Koleksi data dan struct di Golang.",
						VideoURL:    "https://www.youtube.com/watch?v=demo4",
						Duration:    900,
						OrderNumber: 4,
					},
					{
						Title:       "Goroutine & Channel",
						Description: "Concurrency di Golang menggunakan goroutine dan channel.",
						VideoURL:    "https://www.youtube.com/watch?v=demo5",
						Duration:    1200,
						OrderNumber: 5,
					},
				},
			},
		},
	}

	course2 := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  categoryFrontend.ID,
		LevelID:     levelIntermediate.ID,
		Title:       "React.js Fundamental",
		Slug:        "reactjs-fundamental",
		Description: "Kuasai dasar-dasar React.js dari component hingga state management.",
		Price:       200000,
		Duration:    18,
		Status:      "published",
		Sections: []models.CourseSection{
			{
				Title:       "Introduction to React",
				OrderNumber: 1,
				Lessons: []models.Lesson{
					{
						Title:       "Apa itu React?",
						Description: "Mengenal konsep dasar React dan component-based architecture.",
						VideoURL:    "https://www.youtube.com/watch?v=react1",
						Duration:    540,
						OrderNumber: 1,
						IsPreview:   true,
					},
					{
						Title:       "JSX & Rendering",
						Description: "Belajar JSX syntax dan cara rendering element di React.",
						VideoURL:    "https://www.youtube.com/watch?v=react2",
						Duration:    660,
						OrderNumber: 2,
					},
					{
						Title:       "Component & Props",
						Description: "Membuat komponen dan passing data dengan props.",
						VideoURL:    "https://www.youtube.com/watch?v=react3",
						Duration:    780,
						OrderNumber: 3,
					},
				},
			},
			{
				Title:       "State & Lifecycle",
				OrderNumber: 2,
				Lessons: []models.Lesson{
					{
						Title:       "useState Hook",
						Description: "Mengelola state lokal dengan useState hook.",
						VideoURL:    "https://www.youtube.com/watch?v=react4",
						Duration:    840,
						OrderNumber: 4,
						Assignment: &models.Assignment{
							Title:       "Build Counter App",
							Description: "Buat aplikasi counter sederhana menggunakan React useState.",
							DueDate:     timePtr(time.Now().Add(5 * 24 * time.Hour)),
						},
					},
					{
						Title:       "useEffect & API Call",
						Description: "Menggunakan useEffect untuk side effect dan fetch data dari API.",
						VideoURL:    "https://www.youtube.com/watch?v=react5",
						Duration:    960,
						OrderNumber: 5,
					},
				},
			},
		},
	}

	course3 := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  categoryData.ID,
		LevelID:     levelAdvanced.ID,
		Title:       "Machine Learning dengan Python",
		Slug:        "machine-learning-python",
		Description: "Pelajari konsep machine learning dan implementasinya menggunakan Python, NumPy, dan Scikit-learn.",
		Price:       350000,
		Duration:    24,
		Status:      "published",
		Sections: []models.CourseSection{
			{
				Title:       "Dasar Machine Learning",
				OrderNumber: 1,
				Lessons: []models.Lesson{
					{
						Title:       "Pengenalan Machine Learning",
						Description: "Mengapa machine learning penting dan jenis-jenisnya.",
						VideoURL:    "https://www.youtube.com/watch?v=ml1",
						Duration:    720,
						OrderNumber: 1,
						IsPreview:   true,
					},
					{
						Title:       "Data Preprocessing",
						Description: "Membersihkan dan menyiapkan data untuk training model.",
						VideoURL:    "https://www.youtube.com/watch?v=ml2",
						Duration:    1080,
						OrderNumber: 2,
					},
					{
						Title:       "Linear Regression",
						Description: "Memahami dan mengimplementasikan linear regression.",
						VideoURL:    "https://www.youtube.com/watch?v=ml3",
						Duration:    1200,
						OrderNumber: 3,
					},
				},
			},
			{
				Title:       "Classification & Clustering",
				OrderNumber: 2,
				Lessons: []models.Lesson{
					{
						Title:       "Decision Tree & Random Forest",
						Description: "Algoritma klasifikasi berbasis tree.",
						VideoURL:    "https://www.youtube.com/watch?v=ml4",
						Duration:    1440,
						OrderNumber: 4,
					},
					{
						Title:       "K-Means Clustering",
						Description: "Mengelompokkan data tanpa label menggunakan K-Means.",
						VideoURL:    "https://www.youtube.com/watch?v=ml5",
						Duration:    1080,
						OrderNumber: 5,
						Quiz: &models.Quiz{
							Title:        "Kuis Machine Learning",
							PassingScore: 70,
							Questions: []models.QuizQuestion{
								{
									Question:     "Apa itu supervised learning?",
									QuestionType: "multiple_choice",
									Options: []models.QuizOption{
										{OptionText: "Belajar dari data berlabel", IsCorrect: true},
										{OptionText: "Belajar tanpa label", IsCorrect: false},
										{OptionText: "Belajar dari reward", IsCorrect: false},
										{OptionText: "Belajar dari environment", IsCorrect: false},
									},
								},
								{
									Question:     "Algoritma clustering yang paling umum digunakan?",
									QuestionType: "multiple_choice",
									Options: []models.QuizOption{
										{OptionText: "Linear Regression", IsCorrect: false},
										{OptionText: "K-Means", IsCorrect: true},
										{OptionText: "Random Forest", IsCorrect: false},
										{OptionText: "SVM", IsCorrect: false},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	course4 := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  categoryMobile.ID,
		LevelID:     levelBeginner.ID,
		Title:       "Flutter untuk Pemula",
		Slug:        "flutter-untuk-pemula",
		Description: "Membangun aplikasi mobile cross-platform menggunakan Flutter dan Dart dari nol.",
		Price:       175000,
		Duration:    15,
		Status:      "published",
		Sections: []models.CourseSection{
			{
				Title:       "Dart Fundamentals",
				OrderNumber: 1,
				Lessons: []models.Lesson{
					{
						Title:       "Mengenal Dart",
						Description: "Bahasa pemrograman Dart yang digunakan di Flutter.",
						VideoURL:    "https://www.youtube.com/watch?v=flutter1",
						Duration:    600,
						OrderNumber: 1,
						IsPreview:   true,
					},
					{
						Title:       "Widget Basics",
						Description: "Memahami konsep widget di Flutter.",
						VideoURL:    "https://www.youtube.com/watch?v=flutter2",
						Duration:    720,
						OrderNumber: 2,
					},
					{
						Title:       "Layout & Styling",
						Description: "Membuat layout responsif dan styling di Flutter.",
						VideoURL:    "https://www.youtube.com/watch?v=flutter3",
						Duration:    840,
						OrderNumber: 3,
					},
				},
			},
			{
				Title:       "Navigation & State",
				OrderNumber: 2,
				Lessons: []models.Lesson{
					{
						Title:       "Navigator & Routing",
						Description: "Navigasi antar halaman di Flutter.",
						VideoURL:    "https://www.youtube.com/watch?v=flutter4",
						Duration:    660,
						OrderNumber: 4,
					},
					{
						Title:       "State Management",
						Description: "Mengelola state dengan Provider dan setState.",
						VideoURL:    "https://www.youtube.com/watch?v=flutter5",
						Duration:    900,
						OrderNumber: 5,
						Assignment: &models.Assignment{
							Title:       "Build Todo App",
							Description: "Buat aplikasi Todo sederhana dengan fitur tambah, hapus, dan toggle selesai.",
							DueDate:     timePtr(time.Now().Add(10 * 24 * time.Hour)),
						},
					},
				},
			},
		},
	}

	course5 := &models.Course{
		MentorID:    mentor.ID,
		CategoryID:  categoryBackend.ID,
		LevelID:     levelIntermediate.ID,
		Title:       "REST API dengan Go & Gin",
		Slug:        "rest-api-go-gin",
		Description: "Membangun REST API yang production-ready menggunakan Go, Gin, dan GORM.",
		Price:       225000,
		Duration:    20,
		Status:      "draft",
		Sections:    []models.CourseSection{},
	}

	courses := []*models.Course{course1, course2, course3, course4, course5}
	for _, c := range courses {
		var existing models.Course
		err := db.Where("slug = ?", c.Slug).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := db.Create(c).Error; err != nil {
			return err
		}
	}

	log.Println("course content seeded successfully")
	return nil
}

func SeedEnrollments(db *gorm.DB) error {
	var student models.User
	if err := db.Where("email = ?", "student@devacademy.com").First(&student).Error; err != nil {
		return err
	}

	var course1 models.Course
	if err := db.Where("slug = ?", "belajar-golang-dari-dasar").First(&course1).Error; err != nil {
		return err
	}

	var course2 models.Course
	if err := db.Where("slug = ?", "reactjs-fundamental").First(&course2).Error; err != nil {
		return err
	}

	var course3 models.Course
	if err := db.Where("slug = ?", "machine-learning-python").First(&course3).Error; err != nil {
		return err
	}

	var existingEnrollment models.Enrollment
	err := db.Where("user_id = ? AND course_id = ?", student.ID, course1.ID).First(&existingEnrollment).Error
	if err == nil {
		log.Println("enrollments already seeded, skipping")
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	enrolledAt := time.Now().Add(-14 * 24 * time.Hour)

	enrollment1 := models.Enrollment{
		UserID:     student.ID,
		CourseID:   course1.ID,
		EnrolledAt: enrolledAt,
	}
	if err := db.Create(&enrollment1).Error; err != nil {
		return err
	}

	enrollment2 := models.Enrollment{
		UserID:     student.ID,
		CourseID:   course2.ID,
		EnrolledAt: time.Now().Add(-7 * 24 * time.Hour),
	}
	if err := db.Create(&enrollment2).Error; err != nil {
		return err
	}

	enrollment3 := models.Enrollment{
		UserID:     student.ID,
		CourseID:   course3.ID,
		EnrolledAt: time.Now().Add(-3 * 24 * time.Hour),
	}
	if err := db.Create(&enrollment3).Error; err != nil {
		return err
	}

	var course1Lessons []models.Lesson
	if err := db.Joins("JOIN course_sections cs ON lessons.section_id = cs.id").
		Where("cs.course_id = ?", course1.ID).
		Order("cs.order_number ASC, lessons.order_number ASC").
		Find(&course1Lessons).Error; err != nil {
		return err
	}

	lessonProgressData := []models.LessonProgress{
		{EnrollmentID: enrollment1.ID, LessonID: course1Lessons[0].ID, IsCompleted: true, WatchedSecond: 600},
		{EnrollmentID: enrollment1.ID, LessonID: course1Lessons[1].ID, IsCompleted: true, WatchedSecond: 480},
		{EnrollmentID: enrollment1.ID, LessonID: course1Lessons[2].ID, IsCompleted: true, WatchedSecond: 720},
		{EnrollmentID: enrollment1.ID, LessonID: course1Lessons[3].ID, IsCompleted: false, WatchedSecond: 450},
	}

	for i := range lessonProgressData {
		if err := db.Create(&lessonProgressData[i]).Error; err != nil {
			return err
		}
	}

	var course2Lessons []models.Lesson
	if err := db.Joins("JOIN course_sections cs ON lessons.section_id = cs.id").
		Where("cs.course_id = ?", course2.ID).
		Order("cs.order_number ASC, lessons.order_number ASC").
		Find(&course2Lessons).Error; err != nil {
		return err
	}

	lessonProgress2 := []models.LessonProgress{
		{EnrollmentID: enrollment2.ID, LessonID: course2Lessons[0].ID, IsCompleted: true, WatchedSecond: 540},
		{EnrollmentID: enrollment2.ID, LessonID: course2Lessons[1].ID, IsCompleted: false, WatchedSecond: 200},
	}

	for i := range lessonProgress2 {
		if err := db.Create(&lessonProgress2[i]).Error; err != nil {
			return err
		}
	}

	log.Println("enrollments and progress seeded successfully")
	return nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
