package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
	OrderRefer      uuid.UUID `json:"order_id" gorm:"type:uuid;index"`
	Order           Order     `gorm:"foreignKey:OrderRefer"`
	PriceAtPurchase float64   `gorm:"type:decimal(10,2);not null" json:"price_at_purchase"`
	Quantity        int       `gorm:"not null" json:"quantity"`
}

func (orderItem *OrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	orderItem.ID = uuid.New()
	return
}
