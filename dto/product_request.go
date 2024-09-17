package dto

import "github.com/google/uuid"

type RequestProduct struct {
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" validate:"required,gt=0"`
	Stock         int       `json:"stock" validate:"required,gt=0"`
	CategoryRefer uuid.UUID `json:"category_id" validate:"required"`
}
