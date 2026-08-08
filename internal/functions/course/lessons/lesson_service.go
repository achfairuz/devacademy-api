package lessons

import (
	"github.com/google/uuid"

	"context"

	"github.com/iyuz/devacademy-api/internal/models"
)

type LessonService interface {
	Create(ctx context.Context, req *LessonRequest, SectionID string) (*models.Lesson, error)
	GetByID(ctx context.Context, id string) (*models.Lesson, error)
	GetBySectionID(ctx context.Context, sectionID string) ([]*models.Lesson, error)
	Update(ctx context.Context, id string, req *LessonRequest, SectionID string) (*models.Lesson, error)
	Delete(ctx context.Context, id string) error
}

type lessonService struct {
	repo LessonRepository
}

func NewLessonService(repo LessonRepository) LessonService {
	return &lessonService{repo: repo}
}

func (s *lessonService) Create(ctx context.Context, req *LessonRequest, SectionID string) (*models.Lesson, error) {
	sectionID, err := uuid.Parse(SectionID)
	if err != nil {
		return nil, err
	}

	lesson := &models.Lesson{
		SectionID:   sectionID,
		Title:       req.Title,
		Description: req.Description,
		VideoURL:    req.VideoURL,
	}

	if err := s.repo.Create(ctx, lesson); err != nil {
		return nil, err
	}
	return lesson, nil
}

func (s *lessonService) GetByID(ctx context.Context, id string) (*models.Lesson, error) {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	lesson, err := s.repo.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}
	return lesson, nil
}

func (s *lessonService) GetBySectionID(ctx context.Context, sectionID string) ([]*models.Lesson, error) {
	sectionUUID, err := uuid.Parse(sectionID)
	if err != nil {
		return nil, err
	}

	lessons, err := s.repo.GetBySectionID(ctx, sectionUUID)
	if err != nil {
		return nil, err
	}
	return lessons, nil
}

func (s *lessonService) Update(ctx context.Context, id string, req *LessonRequest, SectionID string) (*models.Lesson, error) {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	lesson, err := s.repo.FindByID(ctx, lessonID)
	if err != nil {
		return nil, err
	}

	if SectionID != "" {
		sectionID, err := uuid.Parse(SectionID)
		if err != nil {
			return nil, err
		}
		lesson.SectionID = sectionID
	}
	if req.Title != "" {
		lesson.Title = req.Title
	}
	if req.Description != "" {
		lesson.Description = req.Description
	}
	if req.VideoURL != "" {
		lesson.VideoURL = req.VideoURL
	}

	if err := s.repo.Update(ctx, lesson); err != nil {
		return nil, err
	}
	return lesson, nil
}

func (s *lessonService) Delete(ctx context.Context, id string) error {
	lessonID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, lessonID); err != nil {
		return err
	}
	return nil
}
