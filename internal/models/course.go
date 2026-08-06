package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MentorID    uuid.UUID `gorm:"type:uuid;not null;index"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"size:200;not null"`
	Slug        string    `gorm:"size:200;uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Thumbnail   string    `gorm:"type:text"`
	Price       float64   `gorm:"type:decimal(12,2);not null;default:0"`
	Level       string    `gorm:"size:20"`
	Duration    int
	Status      string `gorm:"size:20;not null;default:draft"`
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Mentor   User            `gorm:"foreignKey:MentorID"`
	Category Category        `gorm:"foreignKey:CategoryID"`
	Sections []CourseSection `gorm:"foreignKey:CourseID"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CourseSection struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"size:200;not null"`
	OrderNumber int       `gorm:"column:order_number;not null;default:1"`

	Course  Course   `gorm:"foreignKey:CourseID"`
	Lessons []Lesson `gorm:"foreignKey:SectionID"`
}

func (cs *CourseSection) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}

type Lesson struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	SectionID   uuid.UUID `gorm:"type:uuid;not null;index"`
	Title       string    `gorm:"size:200;not null"`
	Description string    `gorm:"type:text"`
	VideoURL    string    `gorm:"column:video_url;type:text"`
	Duration    int
	OrderNumber int  `gorm:"column:order_number;not null;default:1"`
	IsPreview   bool `gorm:"column:is_preview;not null;default:false"`

	Section    CourseSection `gorm:"foreignKey:SectionID"`
	Files      []LessonFile  `gorm:"foreignKey:LessonID"`
	Quiz       *Quiz         `gorm:"foreignKey:LessonID"`
	Assignment *Assignment   `gorm:"foreignKey:LessonID"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type LessonFile struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LessonID uuid.UUID `gorm:"type:uuid;not null;index"`
	FileName string    `gorm:"column:file_name;size:255;not null"`
	FileURL  string    `gorm:"column:file_url;type:text;not null"`
	FileSize int       `gorm:"column:file_size"`

	Lesson Lesson `gorm:"foreignKey:LessonID"`
}

func (lf *LessonFile) BeforeCreate(tx *gorm.DB) error {
	if lf.ID == uuid.Nil {
		lf.ID = uuid.New()
	}
	return nil
}
