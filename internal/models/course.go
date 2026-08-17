package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Course struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MentorID    uuid.UUID `gorm:"type:uuid;not null;index" json:"mentor_id"`
	CategoryID  uuid.UUID `gorm:"type:uuid;not null;index" json:"category_id"`
	LevelID     uuid.UUID `gorm:"type:uuid;index" json:"level_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Slug        string    `gorm:"size:200;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Thumbnail   string    `gorm:"type:text" json:"thumbnail"`
	Price       float64   `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	Duration    int       `json:"duration"`
	Status      string    `gorm:"size:20;not null;default:draft" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Mentor   User            `gorm:"foreignKey:MentorID" json:"mentor"`
	Category Category        `gorm:"foreignKey:CategoryID" json:"category"`
	Level    Level           `gorm:"foreignKey:LevelID" json:"level"`
	Sections []CourseSection `gorm:"foreignKey:CourseID" json:"sections"`
}

func (c *Course) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type CourseSection struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CourseID    uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	OrderNumber int       `gorm:"column:order_number;not null;default:1" json:"order_number"`

	Course  Course   `gorm:"foreignKey:CourseID" json:"course"`
	Lessons []Lesson `gorm:"foreignKey:SectionID" json:"lessons"`
}

func (cs *CourseSection) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}

type Lesson struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SectionID   uuid.UUID `gorm:"type:uuid;not null;index" json:"section_id"`
	Title       string    `gorm:"size:200;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	VideoURL    string    `gorm:"column:video_url;type:text" json:"video_url"`
	Duration    int       `json:"duration"`
	OrderNumber int       `gorm:"column:order_number;not null;default:1" json:"order_number"`
	IsPreview   bool      `gorm:"column:is_preview;not null;default:false" json:"is_preview"`

	Section    CourseSection `gorm:"foreignKey:SectionID" json:"section"`
	Files      []LessonFile  `gorm:"foreignKey:LessonID" json:"files"`
	Quiz       *Quiz         `gorm:"foreignKey:LessonID" json:"quiz"`
	Assignment *Assignment   `gorm:"foreignKey:LessonID" json:"assignment"`
}

func (l *Lesson) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type LessonFile struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	LessonID uuid.UUID `gorm:"type:uuid;not null;index" json:"lesson_id"`
	FileName string    `gorm:"column:file_name;size:255;not null" json:"file_name"`
	FileURL  string    `gorm:"column:file_url;type:text;not null" json:"file_url"`
	FileSize int       `gorm:"column:file_size" json:"file_size"`

	Lesson Lesson `gorm:"foreignKey:LessonID" json:"lesson"`
}

func (lf *LessonFile) BeforeCreate(tx *gorm.DB) error {
	if lf.ID == uuid.Nil {
		lf.ID = uuid.New()
	}
	return nil
}
