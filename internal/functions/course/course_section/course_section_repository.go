package coursesection

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type CourseSectionRepository interface {
	Create(ctx context.Context, courseSection *models.CourseSection) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.CourseSection, error)
	FindByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]models.CourseSection, error)
	Update(ctx context.Context, courseSection *models.CourseSection) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type courseSectionRepository struct {
	*repositories.BaseRepository[models.CourseSection]
	db *gorm.DB
}

func NewCourseSectionRepository(db *gorm.DB) CourseSectionRepository {
	return &courseSectionRepository{
		BaseRepository: repositories.NewBaseRepository[models.CourseSection](db),
		db:             db,
	}
}

func (r *courseSectionRepository) FindByCourse(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]models.CourseSection, error) {
	var sections []models.CourseSection
	if err := r.db.WithContext(ctx).
		Where("course_id = ?", courseID).
		Order("order_number ASC").
		Limit(limit).
		Offset(offset).
		Find(&sections).Error; err != nil {
		return nil, err
	}
	return sections, nil
}
