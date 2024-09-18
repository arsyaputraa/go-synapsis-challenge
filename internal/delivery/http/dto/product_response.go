package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
)

type ResponseProduct struct {
	ID          uuid.UUID        `json:"id"` // Use UUID as the primary key
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Price       float64          `json:"price"`
	Stock       int              `json:"stock"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	Category    ResponseCategory `json:"category"`
}

func NewResponseProduct(p *models.Product) ResponseProduct {
	return ResponseProduct{ID: p.ID, Name: p.Name, Description: p.Description, Price: p.Price, Stock: p.Stock, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, Category: NewResponseCategory(&p.Category)}
}
