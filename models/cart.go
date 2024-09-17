package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UserRefer    uuid.UUID     `json:"user_id" gorm:"type:uuid;"`
	User         User    `gorm:"foreignKey:UserRefer"`
}

func (cart *Cart) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	cart.ID = uuid.New()
	return
}