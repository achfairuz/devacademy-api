package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Enrollment struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_enrollments_user_course" json:"user_id"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_enrollments_user_course" json:"course_id"`
	EnrolledAt  time.Time `gorm:"column:enrolled_at;default:now()" json:"enrolled_at"`
	CompletedAt *time.Time `json:"completed_at"`

	User           User             `gorm:"foreignKey:UserID" json:"user"`
	Course         Course           `gorm:"foreignKey:CourseID" json:"course"`
	LessonProgress []LessonProgress `gorm:"foreignKey:EnrollmentID" json:"lesson_progress"`
}

func (e *Enrollment) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type LessonProgress struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	EnrollmentID  uuid.UUID `gorm:"type:uuid;not null;index" json:"enrollment_id"`
	LessonID      uuid.UUID `gorm:"type:uuid;not null;index" json:"lesson_id"`
	IsCompleted   bool      `gorm:"column:is_completed;not null;default:false" json:"is_completed"`
	WatchedSecond int       `gorm:"column:watched_second;not null;default:0" json:"watched_second"`
	UpdatedAt     time.Time `json:"updated_at"`

	Enrollment Enrollment `gorm:"foreignKey:EnrollmentID" json:"enrollment"`
	Lesson     Lesson     `gorm:"foreignKey:LessonID" json:"lesson"`
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
