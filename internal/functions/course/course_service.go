package course

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/utils"
)

var (
	ErrCourseNotFound          = errors.New("course not found")
	ErrSlugTaken               = errors.New("slug already exists")
	ErrLevelRequiredForPublish = errors.New("level_id is required before publishing")
)

type CourseService interface {
	Create(ctx context.Context, req *CreateCourseRequest) (*models.Course, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error)
	GetBySlug(ctx context.Context, slug string) (*models.Course, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Course, error)
	GetByMentor(ctx context.Context, mentorID uuid.UUID, page, pageSize int) ([]models.Course, error)
	GetByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]models.Course, error)
	GetByLevel(ctx context.Context, levelID uuid.UUID, page, pageSize int) ([]models.Course, error)
	GetDetailBySlug(ctx context.Context, slug string) (*CourseDetail, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateCourseRequest) (*models.Course, error)
	UpdateStatus(ctx context.Context, slug string, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type courseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) CourseService {
	return &courseService{repo: repo}
}

func (s *courseService) Create(ctx context.Context, req *CreateCourseRequest) (*models.Course, error) {
	if req.Status == "" {
		req.Status = "draft"
	}

	course := &models.Course{
		MentorID:    req.MentorID,
		CategoryID:  req.CategoryID,
		LevelID:     req.LevelID,
		Title:       req.Title,
		Slug:        utils.Slugify(req.Title),
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Price:       req.Price,
		Duration:    req.Duration,
		Status:      req.Status,
	}

	if req.Status == "published" {
		if err := validateForPublish(course); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, course); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return course, nil
}

func (s *courseService) GetByID(ctx context.Context, id uuid.UUID) (*models.Course, error) {
	course, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

func (s *courseService) GetBySlug(ctx context.Context, slug string) (*models.Course, error) {
	course, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

func (s *courseService) GetDetailBySlug(ctx context.Context, slug string) (*CourseDetail, error) {
	course, err := s.repo.FindDetailBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}
	return toCourseDetail(course), nil
}

func (s *courseService) GetAll(ctx context.Context, page, pageSize int) ([]models.Course, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindAll(ctx, pageSize, (page-1)*pageSize)
}

func (s *courseService) GetByMentor(ctx context.Context, mentorID uuid.UUID, page, pageSize int) ([]models.Course, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindByMentor(ctx, mentorID, pageSize, (page-1)*pageSize)
}

func (s *courseService) GetByCategory(ctx context.Context, categoryID uuid.UUID, page, pageSize int) ([]models.Course, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindByCategory(ctx, categoryID, pageSize, (page-1)*pageSize)
}

func (s *courseService) UpdateStatus(ctx context.Context, slug string, status string) error {
	course, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}
	if status == "published" {
		if err := validateForPublish(course); err != nil {
			return err
		}
	}
	return s.repo.UpdateStatus(ctx, slug, status)
}

func (s *courseService) GetByLevel(ctx context.Context, levelID uuid.UUID, page, pageSize int) ([]models.Course, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindByLevel(ctx, levelID, pageSize, (page-1)*pageSize)
}

func (s *courseService) Update(ctx context.Context, id uuid.UUID, req *UpdateCourseRequest) (*models.Course, error) {
	course, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}

	if req.Title != "" && req.Title != course.Title {
		course.Title = req.Title
		course.Slug = utils.Slugify(req.Title)
	}
	if req.Description != "" {
		course.Description = req.Description
	}
	if req.Thumbnail != "" {
		course.Thumbnail = req.Thumbnail
	}
	if req.Price != 0 {
		course.Price = req.Price
	}
	if req.LevelID != uuid.Nil {
		course.LevelID = req.LevelID
	}
	if req.Duration != 0 {
		course.Duration = req.Duration
	}
	if req.Status != "" {
		course.Status = req.Status
	}

	if course.Status == "published" {
		if err := validateForPublish(course); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, course); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return course, nil
}

func (s *courseService) Delete(ctx context.Context, id uuid.UUID) error {
	course, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if course == nil {
		return ErrCourseNotFound
	}
	return s.repo.Delete(ctx, id)
}

func validateForPublish(course *models.Course) error {
	if course.LevelID == uuid.Nil {
		return ErrLevelRequiredForPublish
	}
	return nil
}
