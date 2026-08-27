package lessonprogress

import (
	"context"

	"github.com/google/uuid"
)

type LessonProgressService interface {
	CountCompletedByUser(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error)
	CountCompletedBySection(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (map[uuid.UUID]int, error)
}

type lessonProgressService struct {
	repo LessonProgressRepository
}

func NewLessonProgressService(repo LessonProgressRepository) LessonProgressService {
	return &lessonProgressService{repo: repo}
}

func (s *lessonProgressService) CountCompletedByUser(ctx context.Context, userID uuid.UUID, courseIDs []uuid.UUID) (map[uuid.UUID]int, error) {
	return s.repo.CountCompletedByUser(ctx, userID, courseIDs)
}

func (s *lessonProgressService) CountCompletedBySection(ctx context.Context, userID uuid.UUID, courseID uuid.UUID) (map[uuid.UUID]int, error) {
	return s.repo.CountCompletedBySection(ctx, userID, courseID)
}
