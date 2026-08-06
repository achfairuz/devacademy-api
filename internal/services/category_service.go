package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/dto"
	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"github.com/iyuz/devacademy-api/internal/utils"
)

var (
	ErrSlugTaken        = errors.New("slug already exists")
	ErrCategoryNotFound = errors.New("category not found")
)

type CategoryService interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequest) (*models.Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Category, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCategoryRequest) (*models.Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name: req.Name,
		Slug: utils.Slugify(req.Name),
		Icon: req.Icon,
	}
	if err := s.repo.Create(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category, nil
}

func (s *categoryService) GetAll(ctx context.Context, page, pageSize int) ([]models.Category, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindAll(ctx, pageSize, (page-1)*pageSize)
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" && req.Name != category.Name {
		category.Name = req.Name
		category.Slug = utils.Slugify(req.Name)
	}
	if req.Icon != "" {
		category.Icon = req.Icon
	}

	if err := s.repo.Update(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	category, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return ErrCategoryNotFound
	}
	return s.repo.Delete(ctx, id)
}
