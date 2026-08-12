package course

import (
	"time"

	"github.com/google/uuid"

	"github.com/iyuz/devacademy-api/internal/models"
)

type CourseDetail struct {
	ID          uuid.UUID        `json:"id"`
	MentorID    uuid.UUID        `json:"mentor_id"`
	CategoryID  uuid.UUID        `json:"category_id"`
	LevelID     uuid.UUID        `json:"level_id"`
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Description string           `json:"description"`
	Thumbnail   string           `json:"thumbnail"`
	Price       float64          `json:"price"`
	Duration    int              `json:"duration"`
	Status      string           `json:"status"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Mentor      *UserSummary     `json:"mentor"`
	Category    *CategorySummary `json:"category"`
	Level       *LevelSummary    `json:"level"`
	Sections    []SectionDetail  `json:"sections"`
}

type UserSummary struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
}

type CategorySummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
	Icon string    `json:"icon"`
}

type LevelSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type SectionDetail struct {
	ID          uuid.UUID      `json:"id"`
	CourseID    uuid.UUID      `json:"course_id"`
	Title       string         `json:"title"`
	OrderNumber int            `json:"order_number"`
	Lessons     []LessonDetail `json:"lessons"`
}

type LessonDetail struct {
	ID          uuid.UUID          `json:"id"`
	SectionID   uuid.UUID          `json:"section_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	VideoURL    string             `json:"video_url"`
	Duration    int                `json:"duration"`
	OrderNumber int                `json:"order_number"`
	IsPreview   bool               `json:"is_preview"`
	Files       []LessonFileDetail `json:"files"`
	Quiz        *QuizDetail        `json:"quiz"`
	Assignment  *AssignmentDetail  `json:"assignment"`
}

type LessonFileDetail struct {
	ID       uuid.UUID `json:"id"`
	FileName string    `json:"file_name"`
	FileURL  string    `json:"file_url"`
	FileSize int       `json:"file_size"`
}

type QuizDetail struct {
	ID           uuid.UUID            `json:"id"`
	Title        string               `json:"title"`
	PassingScore int                  `json:"passing_score"`
	Questions    []QuizQuestionDetail `json:"questions"`
}

type QuizQuestionDetail struct {
	ID           uuid.UUID          `json:"id"`
	Question     string             `json:"question"`
	QuestionType string             `json:"question_type"`
	Options      []QuizOptionDetail `json:"options"`
}

type QuizOptionDetail struct {
	ID         uuid.UUID `json:"id"`
	OptionText string    `json:"option_text"`
}

type AssignmentDetail struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

func toCourseDetail(course *models.Course) *CourseDetail {
	return &CourseDetail{
		ID:          course.ID,
		MentorID:    course.MentorID,
		CategoryID:  course.CategoryID,
		LevelID:     course.LevelID,
		Title:       course.Title,
		Slug:        course.Slug,
		Description: course.Description,
		Thumbnail:   course.Thumbnail,
		Price:       course.Price,
		Duration:    course.Duration,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
		Mentor:      toUserSummary(&course.Mentor),
		Category:    toCategorySummary(&course.Category),
		Level:       toLevelSummary(&course.Level),
		Sections:    toSectionDetails(course.Sections),
	}
}

func toUserSummary(user *models.User) *UserSummary {
	return &UserSummary{
		ID:       user.ID,
		FullName: user.FullName,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Role:     string(user.Role),
	}
}

func toCategorySummary(category *models.Category) *CategorySummary {
	return &CategorySummary{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
		Icon: category.Icon,
	}
}

func toLevelSummary(level *models.Level) *LevelSummary {
	return &LevelSummary{
		ID:   level.ID,
		Name: level.Name,
		Slug: level.Slug,
	}
}

func toSectionDetails(sections []models.CourseSection) []SectionDetail {
	if len(sections) == 0 {
		return []SectionDetail{}
	}
	details := make([]SectionDetail, 0, len(sections))
	for _, section := range sections {
		details = append(details, SectionDetail{
			ID:          section.ID,
			CourseID:    section.CourseID,
			Title:       section.Title,
			OrderNumber: section.OrderNumber,
			Lessons:     toLessonDetails(section.Lessons),
		})
	}
	return details
}

func toLessonDetails(lessons []models.Lesson) []LessonDetail {
	if len(lessons) == 0 {
		return []LessonDetail{}
	}
	details := make([]LessonDetail, 0, len(lessons))
	for _, lesson := range lessons {
		details = append(details, LessonDetail{
			ID:          lesson.ID,
			SectionID:   lesson.SectionID,
			Title:       lesson.Title,
			Description: lesson.Description,
			VideoURL:    lesson.VideoURL,
			Duration:    lesson.Duration,
			OrderNumber: lesson.OrderNumber,
			IsPreview:   lesson.IsPreview,
			Files:       toLessonFileDetails(lesson.Files),
			Quiz:        toQuizDetail(lesson.Quiz),
			Assignment:  toAssignmentDetail(lesson.Assignment),
		})
	}
	return details
}

func toLessonFileDetails(files []models.LessonFile) []LessonFileDetail {
	if len(files) == 0 {
		return []LessonFileDetail{}
	}
	details := make([]LessonFileDetail, 0, len(files))
	for _, file := range files {
		details = append(details, LessonFileDetail{
			ID:       file.ID,
			FileName: file.FileName,
			FileURL:  file.FileURL,
			FileSize: file.FileSize,
		})
	}
	return details
}

func toQuizDetail(quiz *models.Quiz) *QuizDetail {
	if quiz == nil {
		return nil
	}
	return &QuizDetail{
		ID:           quiz.ID,
		Title:        quiz.Title,
		PassingScore: quiz.PassingScore,
		Questions:    toQuizQuestionDetails(quiz.Questions),
	}
}

func toQuizQuestionDetails(questions []models.QuizQuestion) []QuizQuestionDetail {
	if len(questions) == 0 {
		return []QuizQuestionDetail{}
	}
	details := make([]QuizQuestionDetail, 0, len(questions))
	for _, question := range questions {
		details = append(details, QuizQuestionDetail{
			ID:           question.ID,
			Question:     question.Question,
			QuestionType: question.QuestionType,
			Options:      toQuizOptionDetails(question.Options),
		})
	}
	return details
}

func toQuizOptionDetails(options []models.QuizOption) []QuizOptionDetail {
	if len(options) == 0 {
		return []QuizOptionDetail{}
	}
	details := make([]QuizOptionDetail, 0, len(options))
	for _, option := range options {
		details = append(details, QuizOptionDetail{
			ID:         option.ID,
			OptionText: option.OptionText,
		})
	}
	return details
}

func toAssignmentDetail(assignment *models.Assignment) *AssignmentDetail {
	if assignment == nil {
		return nil
	}
	return &AssignmentDetail{
		ID:          assignment.ID,
		Title:       assignment.Title,
		Description: assignment.Description,
		DueDate:     assignment.DueDate,
	}
}
