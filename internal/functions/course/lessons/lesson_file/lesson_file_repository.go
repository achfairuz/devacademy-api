package lessonfile

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type LessonFileRepository interface {
	Create(ctx context.Context, file *models.LessonFile) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.LessonFile, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.LessonFile, error)
	Update(ctx context.Context, file *models.LessonFile) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type lessonFileRepository struct {
	*repositories.BaseRepository[models.LessonFile]
	db *gorm.DB
}

func NewLessonFileRepository(db *gorm.DB) LessonFileRepository {
	return &lessonFileRepository{
		BaseRepository: repositories.NewBaseRepository[models.LessonFile](db),
		db:             db,
	}
}

func (r *lessonFileRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.LessonFile, error) {
	var files []*models.LessonFile
	if err := r.db.WithContext(ctx).Where("lesson_id = ?", lessonID).Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}
