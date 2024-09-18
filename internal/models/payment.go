package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

type PaymentMethod string

const (
	Unpaid PaymentStatus = "unpaid"
	Paid   PaymentStatus = "paid"
	Failed PaymentStatus = "failed"
)

const (
	CC    PaymentMethod = "cc"
	Debit PaymentMethod = "debit"
)

type Payment struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	OrderRefer uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Order      Order     `gorm:"foreignKey:OrderRefer"`
	Status     string    `gorm:"default:unpaid" json:"status"`
	Amount     float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Method     string    `gorm:"not null" json:"method"`
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	payment.ID = uuid.New()
	return
}
