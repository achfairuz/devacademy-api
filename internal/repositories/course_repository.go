package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
)

type CourseRepository interface {
	Create(ctx context.Context, course *models.Course) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Course, error)
	FindBySlug(ctx context.Context, slug string) (*models.Course, error)
	FindByMentor(ctx context.Context, mentorID uuid.UUID, limit, offset int) ([]models.Course, error)
	FindByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) ([]models.Course, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Course, error)
	FindByLevel(ctx context.Context, level string, limit, offset int) ([]models.Course, error)
	Update(ctx context.Context, course *models.Course) error
	UpdateStatus(ctx context.Context, slug string, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
