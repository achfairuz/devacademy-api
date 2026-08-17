package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	CourseID          uuid.UUID `gorm:"type:uuid;not null;index" json:"course_id"`
	CertificateNumber string    `gorm:"column:certificate_number;size:100;uniqueIndex;not null" json:"certificate_number"`
	CertificateURL    string    `gorm:"column:certificate_url;type:text" json:"certificate_url"`
	IssuedAt          time.Time `gorm:"column:issued_at;default:now()" json:"issued_at"`

	User   User   `gorm:"foreignKey:UserID" json:"user"`
	Course Course `gorm:"foreignKey:CourseID" json:"course"`
}

func (c *Certificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
