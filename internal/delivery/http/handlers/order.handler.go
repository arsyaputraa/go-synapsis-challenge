package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
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
	if err := getCartItemsListByCartId(&cartItems, &cart.ID, tx); err != nil {
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

	// Create payment record
	var paymentRequest dto.RequestCreatePayment
	if err := c.BodyParser(&paymentRequest); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Invalid payment request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}
	if err := service.CreatePayment(&order, &paymentRequest, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error creating payment", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Clear the cart
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
	response := dto.NewSuccessResponse(nil, "Order created successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// SERVICE FUNCTION
