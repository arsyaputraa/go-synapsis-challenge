package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	CartRefer    uuid.UUID `json:"cart_id" gorm:"type:uuid;index"`
	Cart         Cart      `gorm:"foreignKey:CartRefer"`
	ProductRefer uuid.UUID `json:"product_id" gorm:"type:uuid;index"`
	Product      Product   `gorm:"foreignKey:ProductRefer"`
	Quantity     int       `gorm:"not null" json:"quantity"`
}

func (cartItem *CartItem) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	cartItem.ID = uuid.New()
	return
}
