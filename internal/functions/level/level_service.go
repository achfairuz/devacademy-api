package level

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/utils"
)

var (
	ErrSlugTaken     = errors.New("slug already exists")
	ErrLevelNotFound = errors.New("level not found")
)

type LevelService interface {
	Create(ctx context.Context, req *CreateLevelRequest) (*models.Level, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Level, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Level, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateLevelRequest) (*models.Level, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type levelService struct {
	repo LevelRepository
}

func NewLevelService(repo LevelRepository) LevelService {
	return &levelService{repo: repo}
}

func (s *levelService) Create(ctx context.Context, req *CreateLevelRequest) (*models.Level, error) {
	level := &models.Level{
		Name: req.Name,
		Slug: utils.Slugify(req.Name),
	}
	if err := s.repo.Create(ctx, level); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return level, nil
}

func (s *levelService) GetByID(ctx context.Context, id uuid.UUID) (*models.Level, error) {
	level, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, ErrLevelNotFound
	}
	return level, nil
}

func (s *levelService) GetAll(ctx context.Context, page, pageSize int) ([]models.Level, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	return s.repo.FindAll(ctx, pageSize, (page-1)*pageSize)
}

func (s *levelService) Update(ctx context.Context, id uuid.UUID, req *UpdateLevelRequest) (*models.Level, error) {
	level, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if level == nil {
		return nil, ErrLevelNotFound
	}

	if req.Name != "" && req.Name != level.Name {
		level.Name = req.Name
		level.Slug = utils.Slugify(req.Name)
	}

	if err := s.repo.Update(ctx, level); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSlugTaken
		}
		return nil, err
	}
	return level, nil
}

func (s *levelService) Delete(ctx context.Context, id uuid.UUID) error {
	level, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if level == nil {
		return ErrLevelNotFound
	}
	return s.repo.Delete(ctx, id)
}
