package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LessonID    uuid.UUID `gorm:"type:uuid;not null;index" json:"lesson_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	DueDate     *time.Time `json:"due_date"`

	Lesson      Lesson                 `gorm:"foreignKey:LessonID" json:"-"`
	Submissions []AssignmentSubmission `gorm:"foreignKey:AssignmentID" json:"-"`
}

func (a *Assignment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type AssignmentSubmission struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AssignmentID uuid.UUID `gorm:"type:uuid;not null;index" json:"assignment_id"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`
	FileURL      string    `gorm:"column:file_url;type:text;not null" json:"file_url"`
	Score        int       `json:"score"`
	Feedback     string    `gorm:"type:text" json:"feedback"`
	SubmittedAt  time.Time `gorm:"column:submitted_at;default:now()" json:"submitted_at"`

	Assignment Assignment `gorm:"foreignKey:AssignmentID" json:"assignment"`
	Student    User        `gorm:"foreignKey:StudentID" json:"student"`
}

func (as *AssignmentSubmission) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}
