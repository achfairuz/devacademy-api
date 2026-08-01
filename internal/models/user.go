package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex;not null"`
	Password  string    `gorm:"size:255;not null"`
	Role      Role      `gorm:"type:varchar(20);not null;default:student;check:chk_users_role,role IN ('admin','instructor','student')"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
