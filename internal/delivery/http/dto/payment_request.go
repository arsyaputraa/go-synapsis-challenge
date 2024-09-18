package dto

import "github.com/google/uuid"

type RequestCreatePayment struct {
	Method string `json:"method" example:"cc" enums:"cc,debit" validate:"required,oneof=cc debit"`
}

type RequestPaymentWebhook struct {
	PaymentID uuid.UUID `json:"payment_id" validate:"required"`
	Status    string    `json:"status" validate:"required"`
	Otp       string    `json:"otp" validate:"required"`
}
