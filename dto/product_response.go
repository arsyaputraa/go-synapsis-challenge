package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
)

type ResponseProduct struct {
	ID          uuid.UUID        `json:"id,omitempty"` // Use UUID as the primary key
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Price       float64          `json:"price,omitempty"`
	Stock       int              `json:"stock,omitempty"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
	Category    ResponseCategory `json:"category,omitempty"`
}

func NewResponseProduct(p *models.Product) ResponseProduct {
	return ResponseProduct{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Stock: p.Stock, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, Category: NewResponseCategory(&p.Category)}
}
