package quiz

import (
	"context"

	"github.com/google/uuid"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"gorm.io/gorm"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz *models.Quiz) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Quiz, error)
	Update(ctx context.Context, quiz *models.Quiz) error
	DeleteQuestions(ctx context.Context, quizID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type quizRepository struct {
	*repositories.BaseRepository[models.Quiz]
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{
		BaseRepository: repositories.NewBaseRepository[models.Quiz](db),
		db:             db,
	}
}

func (r *quizRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	var quiz models.Quiz
	if err := r.db.WithContext(ctx).
		Preload("Questions.Options").
		Where("id = ?", id).
		First(&quiz).Error; err != nil {
		return nil, err
	}
	return &quiz, nil
}

func (r *quizRepository) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Quiz, error) {
	var quizzes []*models.Quiz
	if err := r.db.WithContext(ctx).
		Preload("Questions.Options").
		Where("lesson_id = ?", lessonID).
		Find(&quizzes).Error; err != nil {
		return nil, err
	}
	return quizzes, nil
}

func (r *quizRepository) DeleteQuestions(ctx context.Context, quizID uuid.UUID) error {
	var questionIDs []uuid.UUID
	if err := r.db.WithContext(ctx).
		Model(&models.QuizQuestion{}).
		Where("quiz_id = ?", quizID).
		Pluck("id", &questionIDs).Error; err != nil {
		return err
	}

	if len(questionIDs) > 0 {
		if err := r.db.WithContext(ctx).
			Where("question_id IN ?", questionIDs).
			Delete(&models.QuizOption{}).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).
		Where("quiz_id = ?", quizID).
		Delete(&models.QuizQuestion{}).Error
}

func (r *quizRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("quiz_id = ?", id).Delete(&models.QuizAttempt{}).Error; err != nil {
		return err
	}
	if err := r.DeleteQuestions(ctx, id); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&models.Quiz{}, "id = ?", id).Error
}
