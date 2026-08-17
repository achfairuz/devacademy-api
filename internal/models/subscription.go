package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Price         float64   `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	DurationMonth int       `gorm:"column:duration_month;not null;default:1" json:"duration_month"`
	MaxCourses    int       `gorm:"column:max_courses;not null;default:0" json:"max_courses"`
}

func (sp *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}

type UserSubscription struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	PlanID    uuid.UUID `gorm:"type:uuid;not null;index" json:"plan_id"`
	StartDate time.Time `gorm:"column:start_date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"column:end_date;not null" json:"end_date"`
	Status    string    `gorm:"size:20;not null;default:active" json:"status"`

	User     User             `gorm:"foreignKey:UserID" json:"user"`
	Plan     SubscriptionPlan `gorm:"foreignKey:PlanID" json:"plan"`
	Payments []Payment        `gorm:"foreignKey:UserSubscriptionID" json:"payments"`
}

func (us *UserSubscription) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

type Payment struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserSubscriptionID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_subscription_id"`
	Amount             float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod      string    `gorm:"column:payment_method;size:50" json:"payment_method"`
	PaymentGateway     string    `gorm:"column:payment_gateway;size:50" json:"payment_gateway"`
	TransactionID      string    `gorm:"column:transaction_id;size:100" json:"transaction_id"`
	Status             string    `gorm:"size:20;not null;default:pending" json:"status"`
	PaidAt             *time.Time `json:"paid_at"`

	UserSubscription UserSubscription `gorm:"foreignKey:UserSubscriptionID" json:"user_subscription"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
