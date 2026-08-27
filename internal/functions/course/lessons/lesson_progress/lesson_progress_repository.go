package lessonprogress

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
)

type LessonProgressRepository interface {
	Create(ctx context.Context, progress *models.LessonProgress) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.LessonProgress, error)
	Update(ctx context.Context, progress *models.LessonProgress) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountCompletedByUser(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error)
	CountCompletedBySection(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (map[uuid.UUID]int, error)
}

type lessonProgressRepository struct {
	*repositories.BaseRepository[models.LessonProgress]
	db *gorm.DB
}

func NewLessonProgressRepository(db *gorm.DB) LessonProgressRepository {
	return &lessonProgressRepository{
		BaseRepository: repositories.NewBaseRepository[models.LessonProgress](db),
		db:             db,
	}
}

func (r *lessonProgressRepository) CountCompletedByUser(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error) {
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

func (r *lessonProgressRepository) CountCompletedBySection(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (map[uuid.UUID]int, error) {
	result := make(map[uuid.UUID]int)

	var rows []struct {
		SectionID uuid.UUID
		Total     int
	}
	if err := r.db.WithContext(ctx).
		Model(&models.LessonProgress{}).
		Select("l.section_id, COUNT(*) AS total").
		Joins("JOIN enrollments e ON lesson_progress.enrollment_id = e.id").
		Joins("JOIN lessons l ON lesson_progress.lesson_id = l.id").
		Joins("JOIN course_sections cs ON l.section_id = cs.id").
		Where("e.user_id = ? AND e.course_id = ? AND lesson_progress.is_completed = TRUE", userID, courseID).
		Group("l.section_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.SectionID] = row.Total
	}
	return result, nil
}
