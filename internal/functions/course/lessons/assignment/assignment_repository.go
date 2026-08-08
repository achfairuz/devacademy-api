package assignment

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type AssignmentRepository interface {
	Create(ctx context.Context, assignment *models.Assignment) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Assignment, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Assignment, error)
	Update(ctx context.Context, assignment *models.Assignment) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type assignmentRepository struct {
	*repositories.BaseRepository[models.Assignment]
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{
		BaseRepository: repositories.NewBaseRepository[models.Assignment](db),
		db:             db,
	}
}

func (r *assignmentRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Assignment, error) {
	var assignments []*models.Assignment
	if err := r.db.WithContext(ctx).
		Where("lesson_id = ?", lessonID).
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}
