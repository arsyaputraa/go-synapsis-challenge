package service

import (
	"fmt"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"gorm.io/gorm"
)

func CreatePayment(order *models.Order, paymentRequest *dto.RequestCreatePayment, tx *gorm.DB) error {
	payment := models.Payment{
		OrderRefer: order.ID,
		Status:     string(models.Unpaid),
		Amount:     order.TotalAmount,
		Method:     string(paymentRequest.Method),
	}
	if err := tx.Create(&payment).Error; err != nil {
		return fmt.Errorf("error creating payment: %w", err)
	}
	return nil
}
