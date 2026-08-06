package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_enrollments_user_course"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_enrollments_user_course"`
	EnrolledAt  time.Time `gorm:"column:enrolled_at;default:now()"`
	CompletedAt *time.Time

	User           User             `gorm:"foreignKey:UserID"`
	Course         Course           `gorm:"foreignKey:CourseID"`
	LessonProgress []LessonProgress `gorm:"foreignKey:EnrollmentID"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type LessonProgress struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EnrollmentID  uuid.UUID `gorm:"type:uuid;not null;index"`
	LessonID      uuid.UUID `gorm:"type:uuid;not null;index"`
	IsCompleted   bool      `gorm:"column:is_completed;not null;default:false"`
	WatchedSecond int       `gorm:"column:watched_second;not null;default:0"`
	UpdatedAt     time.Time

	Enrollment Enrollment `gorm:"foreignKey:EnrollmentID"`
	Lesson     Lesson     `gorm:"foreignKey:LessonID"`
}

func (lp *LessonProgress) BeforeCreate(tx *gorm.DB) error {
	if lp.ID == uuid.Nil {
		lp.ID = uuid.New()
	}
	return nil
}

func (lp *LessonProgress) TableName() string {
	return "lesson_progress"
}
