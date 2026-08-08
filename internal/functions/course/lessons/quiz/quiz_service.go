package quiz

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
)

var ErrQuizNotFound = errors.New("quiz not found")

type QuizService interface {
	Create(ctx context.Context, lessonID uuid.UUID, req *QuizRequest) (*models.Quiz, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Quiz, error)
	Update(ctx context.Context, id uuid.UUID, req *QuizRequest) (*models.Quiz, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type quizService struct {
	repo QuizRepository
}

func NewQuizService(repo QuizRepository) QuizService {
	return &quizService{repo: repo}
}

func buildQuestions(req *QuizRequest, quizID uuid.UUID) []models.QuizQuestion {
	questions := make([]models.QuizQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		questionType := q.QuestionType
		if questionType == "" {
			questionType = "multiple_choice"
		}

		question := models.QuizQuestion{
			QuizID:       quizID,
			Question:     q.Question,
			QuestionType: questionType,
		}
		for _, o := range q.Options {
			question.Options = append(question.Options, models.QuizOption{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			})
		}
		questions = append(questions, question)
	}
	return questions
}

func (s *quizService) Create(ctx context.Context, lessonID uuid.UUID, req *QuizRequest) (*models.Quiz, error) {
	quiz := &models.Quiz{
		LessonID:     lessonID,
		Title:        req.Title,
		PassingScore: req.PassingScore,
		Questions:    buildQuestions(req, uuid.Nil),
	}

	if err := s.repo.Create(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) GetByID(ctx context.Context, id uuid.UUID) (*models.Quiz, error) {
	quiz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if quiz == nil {
		return nil, ErrQuizNotFound
	}
	return quiz, nil
}

func (s *quizService) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Quiz, error) {
	return s.repo.GetByLessonID(ctx, lessonID)
}

func (s *quizService) Update(ctx context.Context, id uuid.UUID, req *QuizRequest) (*models.Quiz, error) {
	quiz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if quiz == nil {
		return nil, ErrQuizNotFound
	}

	quiz.Title = req.Title
	quiz.PassingScore = req.PassingScore

	if err := s.repo.DeleteQuestions(ctx, id); err != nil {
		return nil, err
	}
	quiz.Questions = buildQuestions(req, id)

	if err := s.repo.Update(ctx, quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (s *quizService) Delete(ctx context.Context, id uuid.UUID) error {
	quiz, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if quiz == nil {
		return ErrQuizNotFound
	}
	return s.repo.Delete(ctx, id)
}
