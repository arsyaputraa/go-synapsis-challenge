package dto

type RequestCreatePayment struct {
	Method string `json:"method" example:"cc" enums:"cc,debit"`
}
