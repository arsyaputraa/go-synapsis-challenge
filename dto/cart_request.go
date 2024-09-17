package dto

import "github.com/google/uuid"

type RequestAddProductToCart struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}
