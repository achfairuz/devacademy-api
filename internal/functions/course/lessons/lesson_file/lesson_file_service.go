package lessonfile

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
)

var ErrLessonFileNotFound = errors.New("lesson file not found")

type LessonFileService interface {
	Create(ctx context.Context, lessonID uuid.UUID, fileName, fileURL string, fileSize int) (*models.LessonFile, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.LessonFile, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.LessonFile, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type lessonFileService struct {
	repo LessonFileRepository
}

func NewLessonFileService(repo LessonFileRepository) LessonFileService {
	return &lessonFileService{repo: repo}
}

func (s *lessonFileService) Create(ctx context.Context, lessonID uuid.UUID, fileName, fileURL string, fileSize int) (*models.LessonFile, error) {
	file := &models.LessonFile{
		LessonID: lessonID,
		FileName: fileName,
		FileURL:  fileURL,
		FileSize: fileSize,
	}

	if err := s.repo.Create(ctx, file); err != nil {
		return nil, err
	}
	return file, nil
}

func (s *lessonFileService) GetByID(ctx context.Context, id uuid.UUID) (*models.LessonFile, error) {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if file == nil {
		return nil, ErrLessonFileNotFound
	}
	return file, nil
}

func (s *lessonFileService) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.LessonFile, error) {
	return s.repo.GetByLessonID(ctx, lessonID)
}

func (s *lessonFileService) Delete(ctx context.Context, id uuid.UUID) error {
	file, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if file == nil {
		return ErrLessonFileNotFound
	}
	return s.repo.Delete(ctx, id)
}
