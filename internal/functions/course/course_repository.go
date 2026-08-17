package course

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
)

type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Course, error)
	FindBySlug(ctx context.Context, slug string) (*models.Course, error)
	FindByMentor(ctx context.Context, mentorID uuid.UUID, limit, offset int) ([]models.Course, error)
	FindByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]models.Course, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Course, error)
	FindByLevel(ctx context.Context, levelID uuid.UUID, limit, offset int) ([]models.Course, error)
	FindDetailBySlug(ctx context.Context, slug string) (*models.Course, error)
	FindCards(ctx context.Context, limit, offset int) ([]models.Course, error)
	CountEnrollments(ctx context.Context, courseIDs []uuid.UUID) (map[uuid.UUID]int, error)
	CountLessons(ctx context.Context, courseIDs []uuid.UUID) (map[uuid.UUID]int, error)
	CountCompletedLessons(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error)
	Update(ctx context.Context, course *models.Course) error
	UpdateStatus(ctx context.Context, slug string, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type courseRepository struct {
	*repositories.BaseRepository[models.Course]
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) CourseRepository {
	return &courseRepository{
		BaseRepository: repositories.NewBaseRepository[models.Course](db),
		db:             db,
	}
}

func (r *courseRepository) FindBySlug(ctx context.Context, slug string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Where("slug = ?", slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("Sections.Lessons.Files").
		Preload("Sections.Lessons.Quiz.Questions.Options").
		Preload("Sections.Lessons.Assignment").
		Where("id = ?", id).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindCards(ctx context.Context, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Preload("Sections").
		Where("status = ?", "published").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) CountEnrollments(ctx context.Context, courseIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if len(courseIDs) == 0 {
		return result, nil
	}

	var rows []struct {
		CourseID uuid.UUID
		Total    int
	}
	if err := r.db.WithContext(ctx).
		Model(&models.Enrollment{}).
		Select("course_id, COUNT(*) AS total").
		Where("course_id IN ?", courseIDs).
		Group("course_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.CourseID] = row.Total
	}
	return result, nil
}

func (r *courseRepository) CountLessons(ctx context.Context, courseIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if len(courseIDs) == 0 {
		return result, nil
	}

	var rows []struct {
		CourseID uuid.UUID
		Total    int
	}
	if err := r.db.WithContext(ctx).
		Model(&models.Lesson{}).
		Select("cs.course_id, COUNT(*) AS total").
		Joins("JOIN course_sections cs ON lessons.section_id = cs.id").
		Where("cs.course_id IN ?", courseIDs).
		Group("cs.course_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.CourseID] = row.Total
	}
	return result, nil
}

func (r *courseRepository) CountCompletedLessons(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)
	if len(courseIDs) == 0 {
		return result, nil
	}

	var rows []struct {
		CourseID uuid.UUID
		Total    int
	}
	if err := r.db.WithContext(ctx).
		Model(&models.LessonProgress{}).
		Select("cs.course_id, COUNT(*) AS total").
		Joins("JOIN enrollments e ON lesson_progress.enrollment_id = e.id").
		Joins("JOIN lessons l ON lesson_progress.lesson_id = l.id").
		Joins("JOIN course_sections cs ON l.section_id = cs.id").
		Where("e.user_id = ? AND lesson_progress.is_completed = TRUE AND cs.course_id IN ?", userID, courseIDs).
		Group("cs.course_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.CourseID] = row.Total
	}
	return result, nil
}

func (r *courseRepository) FindDetailBySlug(ctx context.Context, slug string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("Sections.Lessons", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_number ASC")
		}).
		Preload("Sections.Lessons.Files").
		Preload("Sections.Lessons.Quiz.Questions.Options").
		Preload("Sections.Lessons.Assignment").
		Where("slug = ?", slug).
		First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByMentor(ctx context.Context, mentorID uuid.UUID, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Where("mentor_id = ?", mentorID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByLevel(ctx context.Context, levelID uuid.UUID, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Preload("Mentor").
		Preload("Category").
		Preload("Level").
		Where("level_id = ?", levelID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) UpdateStatus(ctx context.Context, slug string, status string) error {
	result, err := r.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if result == nil {
		return gorm.ErrRecordNotFound
	}
	result.Status = status
	return r.db.WithContext(ctx).Save(result).Error
}
