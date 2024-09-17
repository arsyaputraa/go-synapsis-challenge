package service

import (
	"fmt"

	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetProducts(products *[]models.Product, query *gorm.DB) error {
	if err := query.Find(products).Error; err != nil {
		return err
	}
	return nil
}

// GetProductByID retrieves a product by its ID
func GetProductByID(product *models.Product, productID uuid.UUID, db *gorm.DB) error {

	if err := db.Preload("Category").First(product, "id = ?", productID).Error; err != nil {
		return err
	}
	return nil
}

func DecrementProductStockByCartItemsQuantity(cartItem models.CartItem, orderID uuid.UUID, tx *gorm.DB) error {
	// Find the product to update stock
	var product models.Product
	if err := tx.First(&product, "id = ?", cartItem.ProductRefer).Error; err != nil {
		return err
	}

	// Check stock availability
	if product.Stock < cartItem.Quantity {
		return fmt.Errorf("insufficient stock for product %s", product.Name)
	}

	// Deduct product stock
	product.Stock -= cartItem.Quantity
	if err := tx.Save(&product).Error; err != nil {
		return fmt.Errorf("error updating product stock: %w", err)
	}

	// Create order item
	orderItem := models.OrderItem{
		OrderRefer:      orderID,
		PriceAtPurchase: product.Price,
		Quantity:        cartItem.Quantity,
	}
	if err := tx.Create(&orderItem).Error; err != nil {
		return fmt.Errorf("error creating order item: %w", err)

	}
	return nil
}

func CreateProduct(newProduct *models.Product, db *gorm.DB) error {
	if err := db.Create(newProduct).Error; err != nil {
		return err
	}
	return nil
}

func UpdateProduct(product *models.Product, db *gorm.DB) error {
	if err := db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func DeleteProduct(product *models.Product, db *gorm.DB) error {
	if err := db.Delete(&product).Error; err != nil {
		return err
	}
	return nil
}
