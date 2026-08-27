package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Quiz struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LessonID     uuid.UUID `gorm:"type:uuid;not null;index" json:"lesson_id"`
	Title        string    `gorm:"size:200;not null" json:"title"`
	PassingScore int       `gorm:"column:passing_score;not null;default:0" json:"passing_score"`

	Lesson    Lesson         `gorm:"foreignKey:LessonID" json:"-"`
	Questions []QuizQuestion `gorm:"foreignKey:QuizID" json:"questions"`
	Attempts  []QuizAttempt  `gorm:"foreignKey:QuizID" json:"-"`
}

func (q *Quiz) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type QuizQuestion struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	QuizID       uuid.UUID `gorm:"type:uuid;not null;index" json:"quiz_id"`
	Question     string    `gorm:"type:text;not null" json:"question"`
	QuestionType string    `gorm:"column:question_type;size:20;not null;default:multiple_choice" json:"question_type"`

	Quiz    Quiz         `gorm:"foreignKey:QuizID" json:"-"`
	Options []QuizOption `gorm:"foreignKey:QuestionID" json:"options"`
}

func (qq *QuizQuestion) BeforeCreate(tx *gorm.DB) error {
	if qq.ID == uuid.Nil {
		qq.ID = uuid.New()
	}
	return nil
}

type QuizOption struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	QuestionID uuid.UUID `gorm:"type:uuid;not null;index" json:"question_id"`
	OptionText string    `gorm:"column:option_text;type:text;not null" json:"option_text"`
	IsCorrect  bool      `gorm:"column:is_correct;not null;default:false" json:"is_correct"`

	Question QuizQuestion `gorm:"foreignKey:QuestionID" json:"-"`
}

func (qo *QuizOption) BeforeCreate(tx *gorm.DB) error {
	if qo.ID == uuid.Nil {
		qo.ID = uuid.New()
	}
	return nil
}

type QuizAttempt struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	QuizID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"quiz_id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Score      int        `gorm:"not null;default:0" json:"score"`
	Status     string     `gorm:"size:20;not null;default:started" json:"status"`
	StartedAt  time.Time  `gorm:"column:started_at;default:now()" json:"started_at"`
	FinishedAt *time.Time `json:"finished_at"`

	Quiz Quiz `gorm:"foreignKey:QuizID" json:"quiz"`
	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (qa *QuizAttempt) BeforeCreate(tx *gorm.DB) error {
	if qa.ID == uuid.Nil {
		qa.ID = uuid.New()
	}
	return nil
}
