package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
)

type ResponseOrder struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Status      string    `json:"status"`
	TotalAmount float64   `json:"total_amount" validate:"required,gt=0"`
}

type ResponseCheckoutOrder struct {
	ID          uuid.UUID `json:"id"`
	Otp         string    `json:"otp"`
	TotalAmount float64   `json:"total_amount" validate:"required,gt=0"`
	PaymentID   uuid.UUID `json:"payment_id"`
}

func NewResponseOrder(p *models.Order) ResponseOrder {
	return ResponseOrder{ID: p.ID, CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt, Status: p.Status, TotalAmount: p.TotalAmount}
}
