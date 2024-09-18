package service

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateOrder(order *models.Order, userID *uuid.UUID, totalAmount float64, tx *gorm.DB) error {
	order.UserRefer = *userID
	order.Status = string(models.Pending)
	order.TotalAmount = totalAmount
	if err := tx.Create(order).Error; err != nil {
		return err
	}
	return nil
}

func GetUserOrders(userID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := database.Database.Db.Where("user_refer = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func GetOrderById(orderID uuid.UUID, db *gorm.DB) (*models.Order, error) {
	var order models.Order
	if err := db.First(&order, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func SaveOrder(order *models.Order, db *gorm.DB) error {
	if err := db.Save(order).Error; err != nil {
		return err
	}
	return nil
}
