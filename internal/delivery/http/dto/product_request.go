package dto

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
)

type RequestProduct struct {
	Name          string    `json:"name" validate:"required"`
	Description   string    `json:"description"`
	Price         float64   `json:"price" validate:"required,gt=0"`
	Stock         int       `json:"stock" validate:"required,gt=0"`
	CategoryRefer uuid.UUID `json:"category_id" validate:"required"`
}

type RequestUpdateProduct struct {
	Name          string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Price         float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock         int       `json:"stock,omitempty" validate:"omitempty,gt=0"`
	CategoryRefer uuid.UUID `json:"category_id,omitempty"`
}

func (rp *RequestProduct) ToModel() models.Product {
	return models.Product{
		Name:          rp.Name,
		Description:   rp.Description,
		Price:         rp.Price,
		Stock:         rp.Stock,
		CategoryRefer: rp.CategoryRefer}
}
