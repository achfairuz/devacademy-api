package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LessonID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Title        string    `gorm:"size:200;not null"`
	PassingScore int       `gorm:"column:passing_score;not null;default:0"`

	Lesson    Lesson         `gorm:"foreignKey:LessonID"`
	Questions []QuizQuestion `gorm:"foreignKey:QuizID"`
	Attempts  []QuizAttempt  `gorm:"foreignKey:QuizID"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type QuizQuestion struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuizID       uuid.UUID `gorm:"type:uuid;not null;index"`
	Question     string    `gorm:"type:text;not null"`
	QuestionType string    `gorm:"column:question_type;size:20;not null;default:multiple_choice"`

	Quiz    Quiz         `gorm:"foreignKey:QuizID"`
	Options []QuizOption `gorm:"foreignKey:QuestionID"`
}

func (qq *QuizQuestion) BeforeCreate(tx *gorm.DB) error {
	if qq.ID == uuid.Nil {
		qq.ID = uuid.New()
	}
	return nil
}

type QuizOption struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null;index"`
	OptionText string    `gorm:"column:option_text;type:text;not null"`
	IsCorrect  bool      `gorm:"column:is_correct;not null;default:false"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID"`
}

func (qo *QuizOption) BeforeCreate(tx *gorm.DB) error {
	if qo.ID == uuid.Nil {
		qo.ID = uuid.New()
	}
	return nil
}

type QuizAttempt struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	QuizID     uuid.UUID `gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Score      int       `gorm:"not null;default:0"`
	Status     string    `gorm:"size:20;not null;default:started"`
	StartedAt  time.Time `gorm:"column:started_at;default:now()"`
	FinishedAt *time.Time

	Quiz Quiz `gorm:"foreignKey:QuizID"`
	User User `gorm:"foreignKey:UserID"`
}

func (qa *QuizAttempt) BeforeCreate(tx *gorm.DB) error {
	if qa.ID == uuid.Nil {
		qa.ID = uuid.New()
	}
	return nil
}
