package service

import (
	"fmt"

	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FindCartByUserId(cart *models.Cart, userID *uuid.UUID, query *gorm.DB) error {
	if err := query.First(cart, "user_refer = ?", userID).Error; err != nil {
		query.Rollback()
		return err
	}
	return nil
}

func ClearCart(cartID uuid.UUID, tx *gorm.DB) error {
	if err := tx.Where("cart_refer = ?", cartID).Delete(&models.CartItem{}).Error; err != nil {
		return fmt.Errorf("Error clearing cart: %w", err)
	}
	return nil
}
