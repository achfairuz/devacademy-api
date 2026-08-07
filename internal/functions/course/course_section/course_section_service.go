package coursesection

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
)

var (
	ErrCourseNotFound        = errors.New("course not found")
	ErrCourseSectionNotFound = errors.New("course section not found")
)

type CourseSectionService interface {
	Create(ctx context.Context, courseID uuid.UUID, req *CourseSectionDto) (*models.CourseSection, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CourseSection, error)
	GetByCourse(ctx context.Context, courseID uuid.UUID, page, pageSize int) ([]models.CourseSection, error)
	Update(ctx context.Context, id uuid.UUID, req *CourseSectionDto) (*models.CourseSection, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type courseSectionService struct {
	repo CourseSectionRepository
}

func NewCourseSectionService(repo CourseSectionRepository) CourseSectionService {
	return &courseSectionService{repo: repo}
}

func (s *courseSectionService) Create(ctx context.Context, courseID uuid.UUID, req *CourseSectionDto) (*models.CourseSection, error) {
	if courseID == uuid.Nil {
		return nil, ErrCourseNotFound
	}

	courseSection := &models.CourseSection{
		CourseID:    courseID,
		Title:       req.Title,
		OrderNumber: req.OrderNumber,
	}

	if err := s.repo.Create(ctx, courseSection); err != nil {
		return nil, err
	}
	return courseSection, nil
}

func (s *courseSectionService) GetByID(ctx context.Context, id uuid.UUID) (*models.CourseSection, error) {
	courseSection, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if courseSection == nil {
		return nil, ErrCourseSectionNotFound
	}
	return courseSection, nil
}

func (s *courseSectionService) GetByCourse(ctx context.Context, courseID uuid.UUID, page, pageSize int) ([]models.CourseSection, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindByCourse(ctx, courseID, pageSize, (page-1)*pageSize)
}

func (s *courseSectionService) Update(ctx context.Context, id uuid.UUID, req *CourseSectionDto) (*models.CourseSection, error) {
	courseSection, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if courseSection == nil {
		return nil, ErrCourseSectionNotFound
	}

	if req.Title != "" {
		courseSection.Title = req.Title
	}
	if req.OrderNumber != 0 {
		courseSection.OrderNumber = req.OrderNumber
	}

	if err := s.repo.Update(ctx, courseSection); err != nil {
		return nil, err
	}
	return courseSection, nil
}

func (s *courseSectionService) Delete(ctx context.Context, id uuid.UUID) error {
	courseSection, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if courseSection == nil {
		return ErrCourseSectionNotFound
	}
	return s.repo.Delete(ctx, id)
}
