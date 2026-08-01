package impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id).Error
}
