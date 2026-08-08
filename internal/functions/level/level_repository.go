package level

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type LevelRepository interface {
	Create(ctx context.Context, level *models.Level) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Level, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.Level, error)
	Update(ctx context.Context, level *models.Level) error
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewLevelRepository(db *gorm.DB) LevelRepository {
	return repositories.NewBaseRepository[models.Level](db)
}
