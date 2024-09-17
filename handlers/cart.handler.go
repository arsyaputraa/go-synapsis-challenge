package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/arsyaputraa/go-synapsis-challenge/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AddToCart godoc
// @Summary Add a product to the shopping cart
// @Description Add a product to the customer's shopping cart
// @Tags cart
// @Accept json
// @Produce json
// @Param cart body dto.RequestAddProductToCart true "Add to cart"
// @Success 200 {object} dto.GeneralResponse "Product added to cart successfully"
// @Failure 400 {object} dto.GeneralResponse "Bad request"
// @Failure 404 {object} dto.GeneralResponse "Product not found"
// @Failure 409 {object} dto.GeneralResponse "Insufficient stock"
// @Router /cart [post]
// @Security BearerAuth
func AddToCart(c *fiber.Ctx) error {
	// Begin a new transaction
	tx := database.Database.Db.Begin()
	if tx.Error != nil {
		response := dto.NewErrorResponse("Error starting transaction", tx.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	userID := c.Locals("userID").(uuid.UUID)

	var addToCartRequest dto.RequestAddProductToCart
	if err := c.BodyParser(&addToCartRequest); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Invalid request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	if err := validate.Struct(&addToCartRequest); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var product models.Product
	if err := tx.First(&product, "id = ?", addToCartRequest.ProductID).Error; err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Product not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Check product stock
	if product.Stock < addToCartRequest.Quantity {
		tx.Rollback()
		response := dto.NewErrorResponse("Insufficient stock", "Not enough stock available for this product")
		return c.Status(fiber.StatusConflict).JSON(response)
	}

	var cart models.Cart
	if err := findOrCreateCartByUserId(&cart, &userID); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error When Finding Cart", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	var cartItem models.CartItem
	if err := tx.Where("cart_refer = ? AND product_refer = ?", cart.ID, product.ID).First(&cartItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			cartItem = models.CartItem{
				CartRefer:    cart.ID,
				ProductRefer: product.ID,
				Quantity:     addToCartRequest.Quantity,
			}
			if err := tx.Create(&cartItem).Error; err != nil {
				tx.Rollback()
				response := dto.NewErrorResponse("Error adding product to cart", err.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(response)
			}
		} else {
			tx.Rollback()
			response := dto.NewErrorResponse("Error finding cart item", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	} else {
		cartItem.Quantity += addToCartRequest.Quantity
		if err := tx.Save(&cartItem).Error; err != nil {
			tx.Rollback()
			response := dto.NewErrorResponse("Error updating cart item", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(response)
		}
	}

	// Update the total cart amount
	if err := updateCartTotalTransaction(&cart, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error updating cart total", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Commit the transaction if everything went well
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error committing transaction", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return a success response
	response := dto.NewSuccessResponse(nil, "Product added to cart successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// Get Cart Items godoc
// @Summary Get Cart Items
// @Description Get a list of cart items
// @Tags cart
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dto.ResponseProduct "cart items data retrieved successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /cart [get]
// @Security BearerAuth
func GetCartItems(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)
	// get cart by user id
	var cart models.Cart
	if err := findOrCreateCartByUserId(&cart, &userID); err != nil {
		response := dto.NewErrorResponse("Error When Finding Cart", err)
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}
	// get cart items by cart id
	query, page, limit := utils.GetPaginatedQuery(&models.CartItem{}, c.Query("page", "1"), c.Query("limit", "10"))
	query = query.Preload("Product")

	var cartItems []models.CartItem
	var totalData int64

	query.Count(&totalData)
	// query cartItems
	if err := query.Find(&cartItems).Error; err != nil {
		response := dto.NewErrorResponse("Error getting cartitems", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	var cartItemDtos []dto.ResponseCartItem
	for _, cartItem := range cartItems {
		productDTO := dto.NewResponseCartItem(&cartItem)
		cartItemDtos = append(cartItemDtos, productDTO)
	}

	// Return the list of transformed products
	paginatedCartItems := dto.ResponsePaginated[dto.ResponseCartItem]{
		Meta: dto.PaginatedMeta{
			Limit: limit,
			Total: int(totalData),
			Page:  page,
		},
		List: cartItemDtos,
	}

	responseCart := dto.NewResponseCart(&cart)
	responseCart.CartItems = paginatedCartItems
	response := dto.NewSuccessResponse(responseCart, "Cart Items Retrieved")
	return c.Status(200).JSON(response)
}

// RemoveCartItems godoc
// @Summary Delete Cart Item
// @Description Delete Cart Item
// @Tags cart
// @Produce json
// @Param id path string true "Cart Items ID"
// @Success 200 {array} dto.ResponseProduct "cart items data deleted successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /cart/{id} [delete]
// @Security BearerAuth
func RemoveCartItems(c *fiber.Ctx) error {
	// Begin a new transaction
	tx := database.Database.Db.Begin()
	if tx.Error != nil {
		response := dto.NewErrorResponse("Error starting transaction", tx.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Get user ID from token
	userID := c.Locals("userID").(uuid.UUID)
	// Get cart items ID from params
	cartItemID := c.Params("id")

	// Find the cart item by its ID
	var cartItem models.CartItem
	if err := tx.Preload("Cart").First(&cartItem, "id = ?", cartItemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			response := dto.NewErrorResponse("Cart item not found", err.Error())
			return c.Status(fiber.StatusNotFound).JSON(response)
		}
		response := dto.NewErrorResponse("Error finding cart item", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Find the cart associated with this cart item and userID
	var cart models.Cart
	if err := tx.First(&cart, "id = ? AND user_refer = ?", cartItem.CartRefer, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {

			response := dto.NewErrorResponse("Cart Not Found", err.Error())
			return c.Status(fiber.StatusNotFound).JSON(response)
		}

		response := dto.NewErrorResponse("Error finding cart", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Delete the cart item
	if err := tx.Delete(&cartItem).Error; err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error removing cart item", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Update the total cart amount
	if err := updateCartTotalTransaction(&cart, tx); err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error updating cart total", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Commit the transaction if everything went well
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		response := dto.NewErrorResponse("Error committing transaction", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return a success response
	response := dto.NewSuccessResponse(nil, "Cart item removed successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// GetCartItem godoc
// @Summary Get Cart Items by ID
// @Description Retrieve a single cart item by its ID
// @Tags cart
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.GeneralResponse "User data retrieved successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /cart/{id} [get]
// @Security BearerAuth
func GetCartItemById(c *fiber.Ctx) error {

	productID := c.Params("id")

	// Fetch the user from the database
	var cartItem models.CartItem
	if err := database.Database.Db.Preload("Product").First(&cartItem, "id = ?", productID).Error; err != nil {
		response := dto.NewErrorResponse("Product not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	// Return the user details
	response := dto.NewSuccessResponse(dto.NewResponseCartItem(&cartItem), "cart item data retrieved successfully")
	return c.Status(fiber.StatusOK).JSON(response)

}

// Updated the function to use the transaction
func updateCartTotalTransaction(cart *models.Cart, tx *gorm.DB) error {
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

// HELPER FUNCTION

func updateCartTotal(cart *models.Cart) error {
	var total float64
	if err := database.Database.Db.Model(&models.CartItem{}).Where("cart_refer = ?", cart.ID).
		Select("COALESCE(SUM(quantity * price), 0)").Joins("JOIN products ON cart_items.product_refer = products.id").Scan(&total).Error; err != nil {
		return err
	}
	cart.TotalAmount = total
	return database.Database.Db.Save(cart).Error
}

func findOrCreateCartByUserId(cart *models.Cart, userID *uuid.UUID) error {

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
