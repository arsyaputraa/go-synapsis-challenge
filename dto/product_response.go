package dto

import (
	"time"

	"github.com/google/uuid"
)

type ResponseProduct struct {
	ID            uuid.UUID `json:"id,omitempty"` // Use UUID as the primary key
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Price         float64   `json:"price,omitempty"`
	Stock         int       `json:"stock,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	Category      ResponseCategory  `json:"category,omitempty"`
}