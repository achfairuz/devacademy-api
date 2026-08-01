package services

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
	"github.com/iyuz/devacademy-api/internal/repositories"
	"github.com/iyuz/devacademy-api/internal/utils"
)

var (
	ErrEmailTaken    = errors.New("email already registered")
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidCreds  = errors.New("invalid email or password")
	ErrUserDuplicate = errors.New("user already exists")
)

type AuthResult struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type UserService interface {
	Register(ctx context.Context, name, email, password string) (*AuthResult, error)
	Login(ctx context.Context, email, password string) (*AuthResult, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.User, error)
	Update(ctx context.Context, id uuid.UUID, name, email string) (*models.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	repo   repositories.UserRepository
	secret string
	expiry int
}

func NewUserService(repo repositories.UserRepository, secret string, expiry int) UserService {
	return &userService{repo: repo, secret: secret, expiry: expiry}
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*AuthResult, error) {
	existing, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{Name: name, Email: email, Password: hash}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(s.secret, s.expiry, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: user}, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*AuthResult, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil || !utils.CheckPassword(user.Password, password) {
		return nil, ErrInvalidCreds
	}

	token, err := utils.GenerateToken(s.secret, s.expiry, user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{Token: token, User: user}, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetAll(ctx context.Context, page, pageSize int) ([]models.User, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	return s.repo.FindAll(ctx, pageSize, offset)
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, name, email string) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if email != "" && email != user.Email {
		dup, err := s.repo.FindByEmail(ctx, email)
		if err != nil {
			return nil, err
		}
		if dup != nil {
			return nil, ErrEmailTaken
		}
		user.Email = email
	}
	if name != "" {
		user.Name = name
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}
