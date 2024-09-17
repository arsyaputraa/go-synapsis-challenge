package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/arsyaputraa/go-synapsis-challenge/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

// Get Products By Category godoc
// @Summary Get Product By Category
// @Description Get a list of products by category or return all products if no category is specified
// @Tags products
// @Produce  json
// @Param category_id query string false "Category ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} dto.ResponseProduct "Products data retrieved successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /product [get]
func GetProductList(c *fiber.Ctx) error {
	// Fetch the user from the database
	categoryID := c.Query("category_id")

	query, page, limit := utils.GetPaginatedQuery(&models.Product{}, c.Query("page", "1"), c.Query("limit", "10"))
	query = query.Preload("Category")

	var products []models.Product
	var totalData int64

	if categoryID != "" {

		if _, err := uuid.Parse(categoryID); err != nil {
			response := dto.NewErrorResponse("Invalid Category ID", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(response)
		}
		query = query.Where("category_refer = ?", categoryID)
	}

	query.Count(&totalData)

	if err := query.Find(&products).Error; err != nil {
		response := dto.NewErrorResponse("Error getting products", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	var productDTOs []dto.ResponseProduct
	for _, product := range products {
		productDTO := dto.NewResponseProduct(&product)
		productDTOs = append(productDTOs, productDTO)
	}

	// Return the list of transformed products

	paginatedResponse := dto.NewPaginatedResponse[dto.ResponseProduct]{
		Total: int(totalData),
		Page:  page,
		Limit: limit,
		List:  productDTOs,
	}
	response := dto.NewSuccessResponse(paginatedResponse, "Products Retrieved")
	return c.Status(200).JSON(response)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Retrieve a single product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.GeneralResponse "User data retrieved successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /product/{id} [get]
func GetProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	// Fetch the user from the database
	var product models.Product
	if err := database.Database.Db.Preload("Category").First(&product, "id = ?", productID).Error; err != nil {
		response := dto.NewErrorResponse("Product not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	// Return the user details
	response := dto.NewSuccessResponse(dto.NewResponseProduct(&product), "User data retrieved successfully")
	return c.JSON(response)

}

// AddProduct godoc
// @Summary Add a new product
// @Description Add a new product to the database
// @Tags admin
// @Accept json
// @Produce json
// @Param product body dto.RequestProduct true "Product data"
// @Success 201 {object} dto.ResponseProduct "Product created successfully"
// @Failure 400 {object} dto.GeneralResponse "Bad request"
// @Router /admin/product [post]
// @Security BearerAuth
func AddProduct(c *fiber.Ctx) error {
	// Parse request body
	var requestProduct dto.RequestProduct
	if err := c.BodyParser(&requestProduct); err != nil {
		response := dto.NewErrorResponse("Invalid request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()

	// Validate request data
	if err := validate.Struct(&requestProduct); err != nil {
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Validate Category ID format
	if _, err := uuid.Parse(requestProduct.CategoryRefer.String()); err != nil {
		response := dto.NewErrorResponse("Invalid Category ID", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	// Create a new product model
	newProduct := requestProduct.ToModel()

	// Save the new product to the database
	if err := database.Database.Db.Create(&newProduct).Error; err != nil {
		response := dto.NewErrorResponse("Error creating product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return the created product as a response
	response := dto.NewSuccessResponse(dto.NewResponseProduct(&newProduct), "Product created successfully")
	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product's details. Only accessible by users with the 'admin' role.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body dto.RequestUpdateProduct true "Product data"
// @Success 200 {object} dto.ResponseProduct "Product updated successfully"
// @Failure 400 {object} dto.GeneralResponse "Invalid request"
// @Failure 404 {object} dto.GeneralResponse "Product not found"
// @Failure 403 {object} dto.GeneralResponse "Forbidden. Only admin can access this endpoint."
// @Router /admin/product/{id} [patch]
// @Security BearerAuth
func UpdateProduct(c *fiber.Ctx) error {
	productID := c.Params("id")

	var updateProduct dto.RequestUpdateProduct
	if err := c.BodyParser(&updateProduct); err != nil {
		response := dto.NewErrorResponse("Invalid request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	if err := validate.Struct(&updateProduct); err != nil {
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var product models.Product
	if err := database.Database.Db.First(&product, "id = ?", productID).Error; err != nil {
		response := dto.NewErrorResponse("Product not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err := copier.CopyWithOption(&product, &updateProduct, copier.Option{IgnoreEmpty: true}); err != nil {
		response := dto.NewErrorResponse("Error updating product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if err := database.Database.Db.Save(&product).Error; err != nil {
		response := dto.NewErrorResponse("Error updating product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := dto.NewSuccessResponse(dto.NewResponseProduct(&product), "Product updated successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID. Only accessible by users with the 'admin' role.
// @Tags admin
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.GeneralResponse "Product deleted successfully"
// @Failure 404 {object} dto.GeneralResponse "Product not found"
// @Router /admin/product/{id} [delete]
// @Security BearerAuth
func DeleteProduct(c *fiber.Ctx) error {
	// Get the product ID from the URL path
	productID := c.Params("id")

	// Fetch the product from the database
	var product models.Product
	if err := database.Database.Db.First(&product, "id = ?", productID).Error; err != nil {
		// If the product is not found, return a 404 response
		response := dto.NewErrorResponse("Product not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	// Delete the product
	if err := database.Database.Db.Delete(&product).Error; err != nil {
		response := dto.NewErrorResponse("Error deleting product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	// Return a success response
	response := dto.NewSuccessResponse(nil, "Product deleted successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}
