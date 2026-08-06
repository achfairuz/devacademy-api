package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID            uuid.UUID `gorm:"type:uuid;not null;index"`
	CourseID          uuid.UUID `gorm:"type:uuid;not null;index"`
	CertificateNumber string    `gorm:"column:certificate_number;size:100;uniqueIndex;not null"`
	CertificateURL    string    `gorm:"column:certificate_url;type:text"`
	IssuedAt          time.Time `gorm:"column:issued_at;default:now()"`

	User   User   `gorm:"foreignKey:UserID"`
	Course Course `gorm:"foreignKey:CourseID"`
}

func (c *Certificate) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
