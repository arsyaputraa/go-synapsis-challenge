package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/arsyaputraa/go-synapsis-challenge/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		productDTO := product.ToDto()
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
	response := dto.NewSuccessResponse(product.ToDto(), "User data retrieved successfully")
	return c.JSON(response)

}
