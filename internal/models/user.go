package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	FullName      string         `gorm:"column:full_name;size:100;not null;default:''" json:"full_name"`
	Username      string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email         string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"size:255;not null" json:"-"`
	Phone         string         `gorm:"size:30" json:"phone"`
	Avatar        string         `gorm:"type:text" json:"avatar"`
	Role          Role           `gorm:"type:varchar(20);not null;default:student;check:chk_users_role,role IN ('admin','mentor','student')" json:"role"`
	Status        bool           `gorm:"not null;default:true" json:"status"`
	EmailVerified bool           `gorm:"column:email_verified;not null;default:false" json:"email_verified"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
