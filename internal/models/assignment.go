package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Assignment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LessonID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"size:200;not null"`
	Description string    `gorm:"type:text"`
	DueDate     *time.Time

	Lesson      Lesson                 `gorm:"foreignKey:LessonID"`
	Submissions []AssignmentSubmission `gorm:"foreignKey:AssignmentID"`
}

func (a *Assignment) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type AssignmentSubmission struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AssignmentID uuid.UUID `gorm:"type:uuid;not null;index"`
	StudentID    uuid.UUID `gorm:"type:uuid;not null;index"`
	FileURL      string    `gorm:"column:file_url;type:text;not null"`
	Score        int
	Feedback     string    `gorm:"type:text"`
	SubmittedAt  time.Time `gorm:"column:submitted_at;default:now()"`

	Assignment Assignment `gorm:"foreignKey:AssignmentID"`
	Student    User       `gorm:"foreignKey:StudentID"`
}

func (as *AssignmentSubmission) BeforeCreate(tx *gorm.DB) error {
	if as.ID == uuid.Nil {
		as.ID = uuid.New()
	}
	return nil
}
