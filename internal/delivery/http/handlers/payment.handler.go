package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// PaymentWebhook godoc
// @Summary Mocking Webhook for payment gateway to update payment and order status
// @Description Webhook endpoint to mock update payment status upon completion/failed payment
// @Tags payment
// @Accept json
// @Produce json
// @Param paymentId query string true "Payment ID"
// @Param status query string true "Payment status" Enums(paid, unpaid, failed)
// @Param otp query string true "Payment OTP -> Get from checking out order"
// @Success 200 {object} dto.GeneralResponse "Webhook processed successfully"
// @Failure 400 {object} dto.GeneralResponse "Bad request"
// @Failure 404 {object} dto.GeneralResponse "Payment not found"
// @Failure 500 {object} dto.GeneralResponse "Internal server error"
// @Router /webhook/payment [get]
func PaymentWebhook(c *fiber.Ctx) error {
	paymentID := c.Query("paymentId")
	paymentStatus := c.Query("status")
	paymentOtp := c.Query("otp")

	paymentUUID, err := utils.CheckUUID(paymentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid Payment ID", err))
	}

	var webhookRequest = dto.RequestPaymentWebhook{PaymentID: *paymentUUID, Status: paymentStatus, Otp: paymentOtp}

	// Start a transaction
	tx := database.Database.Db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error starting transaction", tx.Error.Error()))
	}

	payment, err := service.AuthenticatePayment(webhookRequest.PaymentID, webhookRequest.Otp, tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse("Unauthorized Payment", err.Error()))
	}
	// Update payment status
	payment.Status = webhookRequest.Status
	if err := service.SavePayment(payment, tx); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error updating payment status", err.Error()))
	}

	// If payment is successful, update the order status

	if webhookRequest.Status == "paid" {
		order, err := service.GetOrderById(payment.OrderRefer, tx)
		if err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("Order not found", err.Error()))
		}
		order.Status = string(models.PaidOrder)

		if err := service.SaveOrder(order, tx); err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error updating order status", err.Error()))
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error committing transaction", err.Error()))
	}

	// Return a success response
	response := dto.NewSuccessResponse(nil, "Webhook processed successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}
