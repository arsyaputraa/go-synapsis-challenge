package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
	Name          string    `gorm:"type:varchar(100);not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Stock         int       `gorm:"not null" json:"stock"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	CategoryRefer uuid.UUID `json:"category_id" gorm:"type:uuid;"`
	Category      Category  `gorm:"foreignKey:CategoryRefer"`
}

// BeforeCreate hook will be triggered before inserting a new record to the database
func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	product.ID = uuid.New()
	return
}
