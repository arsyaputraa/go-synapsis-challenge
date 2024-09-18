package service

import (
	"fmt"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
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
		return fmt.Errorf("error clearing cart: %w", err)
	}
	return nil
}

func FindOrCreateCartByUserId(cart *models.Cart, userID *uuid.UUID) error {

	if err := database.Database.Db.Where("user_refer = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			newCart := models.Cart{
				UserRefer: *userID,
			}
			if err := database.Database.Db.Create(&newCart).Error; err != nil {
				return err
			}
			*cart = newCart
		} else {
			return err
		}
	}
	return nil
}

func GetCartItemsListByCartId(cartItems *[]models.CartItem, cartID *uuid.UUID, query *gorm.DB) error {

	if err := query.Where("cart_refer = ?", *cartID).Preload("Product").Find(&cartItems).Error; err != nil {
		return err
	}

	return nil
}

// SERVICE FUNCTION
func UpdateCartTotalTransaction(cart *models.Cart, tx *gorm.DB) error {
	var total float64
	if err := tx.Model(&models.CartItem{}).
		Where("cart_refer = ?", cart.ID).
		Select("COALESCE(SUM(quantity * price), 0)"). // Use COALESCE to handle NULL
		Joins("JOIN products ON cart_items.product_refer = products.id").
		Scan(&total).Error; err != nil {
		return err
	}

	cart.TotalAmount = total
	return tx.Save(cart).Error
}
