package service

import (
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
