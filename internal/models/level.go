package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Level struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"size:50;not null"`
	Slug      string    `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt time.Time
}

func (l *Level) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
