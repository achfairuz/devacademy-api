package lessons

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type LessonRepository interface {
	Create(ctx context.Context, lesson *models.Lesson) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Lesson, error)
	GetBySectionID(ctx context.Context, sectionID uuid.UUID) ([]*models.Lesson, error)
	Update(ctx context.Context, lesson *models.Lesson) error
	Delete(ctx context.Context, uuid uuid.UUID) error
}

type lessonRepository struct {
	*repositories.BaseRepository[models.Lesson]
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &lessonRepository{
		BaseRepository: repositories.NewBaseRepository[models.Lesson](db),
		db:             db,
	}
}

func (r *lessonRepository) GetBySectionID(ctx context.Context, sectionID uuid.UUID) ([]*models.Lesson, error) {
	var lessons []*models.Lesson
	if err := r.db.WithContext(ctx).Where("section_id = ?", sectionID).Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}
