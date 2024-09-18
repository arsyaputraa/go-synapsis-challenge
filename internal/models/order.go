package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	PaidOrder OrderStatus = "paid"
	Completed OrderStatus = "completed"
	Canceled  OrderStatus = "canceled"
)

type Order struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UserRefer   uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	User        User      `gorm:"foreignKey:UserRefer"`
	Status      string    `gorm:"default:pending" json:"status"`
	TotalAmount float64   `gorm:"type:decimal(10,2);not null" json:"total_amount"`
}

func (order *Order) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	order.ID = uuid.New()
	return
}
