package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
    ID          uuid.UUID `gorm:"type:uuid;primaryKey"` // Use UUID as the primary key
    Name        string    `gorm:"type:varchar(100);not null"`
    Description string    `gorm:"type:text"`
    CreatedAt   time.Time `gorm:"autoCreateTime"`
    UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

// BeforeCreate hook will be triggered before inserting a new record to the database
func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
    // Generate a new UUID and assign it to the ID field
    category.ID = uuid.New()
    return
}