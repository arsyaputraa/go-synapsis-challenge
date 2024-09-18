package service

import (
	"errors"
	"fmt"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrPaymentNotFound          = errors.New("payment not found")
	ErrInvalidPaymentCredential = errors.New("unauthorized payment")
)

func CreatePayment(order *models.Order, payment *models.Payment, paymentRequest *dto.RequestCreatePayment, otp *string, tx *gorm.DB) error {

	secretCode := utils.GenerateRandomCode(8)

	codeHash, err := bcrypt.GenerateFromPassword([]byte(secretCode), 14)
	if err != nil {
		return err
	}
	*payment = models.Payment{
		OrderRefer: order.ID,
		Status:     string(models.Unpaid),
		Amount:     order.TotalAmount,
		Method:     string(paymentRequest.Method),
		Otp:        string(codeHash),
	}
	if err := tx.Create(payment).Error; err != nil {
		return fmt.Errorf("error creating payment: %w", err)
	}
	*otp = secretCode
	return nil
}

func AuthenticatePayment(id uuid.UUID, otp string, tx *gorm.DB) (*models.Payment, error) {
	payment, err := FindPaymentByID(id, tx)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(payment.Otp), []byte(otp)); err != nil {
		return nil, ErrInvalidPaymentCredential
	}
	return payment, nil
}

func FindPaymentByID(id uuid.UUID, tx *gorm.DB) (*models.Payment, error) {
	var payment models.Payment

	if err := tx.Where("id = ?", id).First(&payment).Error; err != nil {

		return nil, ErrPaymentNotFound
	}
	return &payment, nil
}

func SavePayment(payment *models.Payment, tx *gorm.DB) error {
	if err := tx.Save(&payment).Error; err != nil {
		return err
	}
	return nil
}

// func FindPaymentByID()(*models.Payment)
