package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context, limit, offset int) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
