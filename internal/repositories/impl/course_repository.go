package impl

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
)

type courseRepository struct {
	*repositories.BaseRepository[models.Course]
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) repositories.CourseRepository {
	return &courseRepository{
		BaseRepository: repositories.NewBaseRepository[models.Course](db),
		db:             db,
	}
}

func (r *courseRepository) FindBySlug(ctx context.Context, slug string) (*models.Course, error) {
	var course models.Course
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&course).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &course, nil
}

func (r *courseRepository) FindByMentor(ctx context.Context, mentorID uuid.UUID, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
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
		Where("category_id = ?", categoryID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&courses).Error; err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *courseRepository) FindByLevel(ctx context.Context, level string, limit, offset int) ([]models.Course, error) {
	var courses []models.Course
	if err := r.db.WithContext(ctx).
		Where("level = ?", level).
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
