package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string    `gorm:"size:100;not null"`
	Price         float64   `gorm:"type:decimal(12,2);not null;default:0"`
	DurationMonth int       `gorm:"column:duration_month;not null;default:1"`
	MaxCourses    int       `gorm:"column:max_courses;not null;default:0"`
}

func (sp *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

type UserSubscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null;index"`
	StartDate time.Time `gorm:"column:start_date;not null"`
	EndDate   time.Time `gorm:"column:end_date;not null"`
	Status    string    `gorm:"size:20;not null;default:active"`

	User     User             `gorm:"foreignKey:UserID"`
	Plan     SubscriptionPlan `gorm:"foreignKey:PlanID"`
	Payments []Payment        `gorm:"foreignKey:UserSubscriptionID"`
}

func (us *UserSubscription) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

type Payment struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserSubscriptionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Amount             float64   `gorm:"type:decimal(12,2);not null"`
	PaymentMethod      string    `gorm:"column:payment_method;size:50"`
	PaymentGateway     string    `gorm:"column:payment_gateway;size:50"`
	TransactionID      string    `gorm:"column:transaction_id;size:100"`
	Status             string    `gorm:"size:20;not null;default:pending"`
	PaidAt             *time.Time

	UserSubscription UserSubscription `gorm:"foreignKey:UserSubscriptionID"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
