package assignment

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
)

var ErrAssignmentNotFound = errors.New("assignment not found")

type AssignmentService interface {
	Create(ctx context.Context, lessonID uuid.UUID, req *AssignmentRequest) (*models.Assignment, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Assignment, error)
	GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Assignment, error)
	Update(ctx context.Context, id uuid.UUID, req *AssignmentRequest) (*models.Assignment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type assignmentService struct {
	repo AssignmentRepository
}

func NewAssignmentService(repo AssignmentRepository) AssignmentService {
	return &assignmentService{repo: repo}
}

func (s *assignmentService) Create(ctx context.Context, lessonID uuid.UUID, req *AssignmentRequest) (*models.Assignment, error) {
	assignment := &models.Assignment{
		LessonID:    lessonID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Create(ctx, assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *assignmentService) GetByID(ctx context.Context, id uuid.UUID) (*models.Assignment, error) {
	assignment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, ErrAssignmentNotFound
	}
	return assignment, nil
}

func (s *assignmentService) GetByLessonID(ctx context.Context, lessonID uuid.UUID) ([]*models.Assignment, error) {
	return s.repo.GetByLessonID(ctx, lessonID)
}

func (s *assignmentService) Update(ctx context.Context, id uuid.UUID, req *AssignmentRequest) (*models.Assignment, error) {
	assignment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, ErrAssignmentNotFound
	}

	if req.Title != "" {
		assignment.Title = req.Title
	}
	if req.Description != "" {
		assignment.Description = req.Description
	}
	if req.DueDate != nil {
		assignment.DueDate = req.DueDate
	}

	if err := s.repo.Update(ctx, assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func (s *assignmentService) Delete(ctx context.Context, id uuid.UUID) error {
	assignment, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return ErrAssignmentNotFound
	}
	return s.repo.Delete(ctx, id)
}
