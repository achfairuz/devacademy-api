package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return NewBaseRepository[models.Category](db)
}
