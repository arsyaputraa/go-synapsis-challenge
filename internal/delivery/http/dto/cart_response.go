package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
)

type ResponseCart struct {
	ID          uuid.UUID                           `json:"id"` // Use UUID as the primary key
	UpdatedAt   time.Time                           `json:"updated_at"`
	TotalAmount float64                             `json:"total_amount"`
	CartItems   ResponsePaginated[ResponseCartItem] `json:"cart_items"`
}

func NewResponseCart(c *models.Cart) ResponseCart {
	return ResponseCart{ID: c.ID, UpdatedAt: c.UpdatedAt, TotalAmount: c.TotalAmount, CartItems: ResponsePaginated[ResponseCartItem]{}}
}

type ResponseCartItem struct {
	ID             uuid.UUID `json:"id"` // Use UUID as the primary key
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Quantity       int       `json:"quantity"`
	ProductID      uuid.UUID `json:"product_id"` // Use UUID as the primary key
	Name           string    `json:"product_name"`
	Description    string    `json:"product_description"`
	Price          float64   `json:"product_price"`
	TotalItemPrice float64   `json:"total_item_price"`
}

func NewResponseCartItem(i *models.CartItem) ResponseCartItem {

	totalItemPrice := i.Product.Price * float64(i.Quantity)
	return ResponseCartItem{ID: i.ID, UpdatedAt: i.UpdatedAt, Quantity: i.Quantity, CreatedAt: i.CreatedAt, ProductID: i.Product.ID, Name: i.Product.Name, Description: i.Product.Description, Price: i.Product.Price, TotalItemPrice: totalItemPrice}
}
