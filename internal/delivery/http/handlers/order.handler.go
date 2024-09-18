package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Checkout godoc
// @Summary Checkout cart and create an order
// @Description Checkout cart and create an order
// @Tags order
// @Accept json
// @Produce json
// @Param payment body dto.RequestCreatePayment true "Payment details"
// @Success 200 {object} dto.GeneralResponse "Order created successfully"
// @Failure 400 {object} dto.GeneralResponse "Bad request"
// @Failure 404 {object} dto.GeneralResponse "Cart not found or empty"
// @Failure 409 {object} dto.GeneralResponse "Insufficient stock"
// @Router /order/checkout [post]
// @Security BearerAuth
func CheckoutOrder(c *fiber.Ctx) error {
	// Begin a new transaction
	tx := database.Database.Db.Begin()
	if tx.Error != nil {
		response := dto.NewErrorResponse("Error starting transaction", tx.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Get user ID from token
	userID := c.Locals("userID").(uuid.UUID)

	// Find the user's cart
	var cart models.Cart
	if err := service.FindCartByUserId(&cart, &userID, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Cart not found", "No cart found for this user.")
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// query cartItems
	var cartItems []models.CartItem
	if err := service.GetCartItemsListByCartId(&cartItems, &cart.ID, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error getting cart items", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	// Check if the cart is empty
	if len(cartItems) == 0 || cartItems == nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Cart is empty", "Cannot checkout an empty cart.")
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// update cart total
	if err := service.UpdateCartTotalTransaction(&cart, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error updating cart total", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Create a new order

	var order models.Order
	if err := service.CreateOrder(&order, &userID, cart.TotalAmount, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error creating order", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Process each cart item
	for _, cartItem := range cartItems {
		if err := service.DecrementProductStockByCartItemsQuantity(cartItem, order.ID, tx); err != nil {
			tx.Rollback()
			response := dto.NewErrorResponse("Error processing cart items", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	validate := validator.New()
	// Create payment record
	var paymentRequest dto.RequestCreatePayment
	var otp string
	var payment models.Payment
	if err := c.BodyParser(&paymentRequest); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Invalid payment request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := validate.Struct(&paymentRequest); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	if err := service.CreatePayment(&order, &payment, &paymentRequest, &otp, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error creating payment", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if err := service.ClearCart(cart.ID, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error clearing cart", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Commit the transaction if everything went well
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error committing transaction", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return a success response
	response := dto.NewSuccessResponse(dto.ResponseCheckoutOrder{ID: order.ID, Otp: otp, TotalAmount: order.TotalAmount, PaymentID: payment.ID}, "Order created successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetUserOrders godoc
// @Summary Get current user's orders
// @Description Get the list of orders for the authenticated user
// @Tags order
// @Produce json
// @Success 200 {array} dto.ResponseOrder "List of user orders"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /order [get]
// @Security BearerAuth
func GetUserOrders(c *fiber.Ctx) error {
	// Get user ID from the request context
	userID := c.Locals("userID").(uuid.UUID)
	// Fetch orders using the service layer
	orders, err := service.GetUserOrders(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error retrieving orders", err.Error()))
	}
	var orderResponses []dto.ResponseOrder
	for _, order := range orders {
		orderResponses = append(orderResponses, dto.NewResponseOrder(&order))
	}
	return c.JSON(dto.NewSuccessResponse(orderResponses, "Orders retrieved successfully"))
}
