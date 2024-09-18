package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

// Get Categories godoc
// @Summary Get Categories
// @Description Get a list of Categories
// @Tags category
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Success 200 {array} dto.ResponseProduct "Categories data retrieved successfully"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /category [get]
func GetCategories(c *fiber.Ctx) error {
	query, page, limit := utils.GetPaginatedQuery(&models.Category{}, c.Query("page", "1"), c.Query("limit", "20"))

	var totalData int64

	query.Count(&totalData)
	var categories []models.Category
	if err := service.GetCategories(&categories, query); err != nil {
		response := dto.NewErrorResponse("Error getting products", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}
	var categoryDTOs []dto.ResponseCategory
	for _, category := range categories {
		categoryDTO := dto.NewResponseCategory(&category)
		categoryDTOs = append(categoryDTOs, categoryDTO)
	}

	// Return the list of transformed products

	paginatedResponse := dto.ResponsePaginated[dto.ResponseCategory]{
		Meta: dto.PaginatedMeta{
			Limit: limit,
			Total: int(totalData),
			Page:  page,
		},
		List: categoryDTOs,
	}
	response := dto.NewSuccessResponse(paginatedResponse, "CAtegories Retrieved")
	return c.Status(200).JSON(response)
}

// AddCategory godoc
// @Summary Add a new Category
// @Description Add a new Category to the database
// @Tags admin
// @Accept json
// @Produce json
// @Param Category body dto.RequestCategory true "Category data"
// @Success 201 {object} dto.GeneralResponse "Category created successfully"
// @Failure 400 {object} dto.GeneralResponse "Bad request"
// @Router /admin/category [post]
// @Security BearerAuth
func AddCategory(c *fiber.Ctx) error {
	var requestCategory dto.RequestCategory
	if err := c.BodyParser(&requestCategory); err != nil {
		response := dto.NewErrorResponse("Invalid request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()

	if err := validate.Struct(&requestCategory); err != nil {
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	newCategory := requestCategory.ToModel()

	if err := service.CreateCategory(&newCategory, database.Database.Db); err != nil {
		response := dto.NewErrorResponse("Error creating product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := dto.NewSuccessResponse(dto.NewResponseCategory(&newCategory), "Category created successfully")
	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateCategory godoc
// @Summary Update a Category
// @Description Update a Category's. Only accessible by users with the 'admin' role.
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param Category body dto.RequestUpdateCategory true "Category data"
// @Success 200 {object} dto.GeneralResponse "Category updated successfully"
// @Failure 400 {object} dto.GeneralResponse "Invalid request"
// @Failure 404 {object} dto.GeneralResponse "Product not found"
// @Failure 403 {object} dto.GeneralResponse "Forbidden. Only admin can access this endpoint."
// @Router /admin/category/{id} [patch]
// @Security BearerAuth
func UpdateCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	db := database.Database.Db
	categoryUUID, err := utils.CheckUUID(categoryID)
	if err != nil {
		response := dto.NewErrorResponse("Invalid Category ID", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var updateCategory dto.RequestUpdateCategory
	if err := c.BodyParser(&updateCategory); err != nil {
		response := dto.NewErrorResponse("Invalid request", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	validate := validator.New()
	if err := validate.Struct(&updateCategory); err != nil {
		response := dto.NewErrorResponse("Validation error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var category models.Category
	if err := service.GetCategoryByID(&category, *categoryUUID, db); err != nil {
		response := dto.NewErrorResponse("Category not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err := copier.CopyWithOption(&category, &updateCategory, copier.Option{IgnoreEmpty: true}); err != nil {
		response := dto.NewErrorResponse("Error updating category", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if err := service.UpdateCategory(&category, db); err != nil {
		response := dto.NewErrorResponse("Error updating product", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := dto.NewSuccessResponse(dto.NewResponseCategory(&category), "Category updated successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteCategory godoc
// @Summary Delete a Category
// @Description Delete a Category by its ID. Only accessible by users with the 'admin' role.
// @Tags admin
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} dto.GeneralResponse "Category deleted successfully"
// @Failure 404 {object} dto.GeneralResponse "Category not found"
// @Router /admin/category/{id} [delete]
// @Security BearerAuth
func DeleteCategory(c *fiber.Ctx) error {
	categoryID := c.Params("id")
	db := database.Database.Db
	categoryUUID, err := utils.CheckUUID(categoryID)
	if err != nil {
		response := dto.NewErrorResponse("Invalid Category ID", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(response)
	}

	var category models.Category
	if err := service.GetCategoryByID(&category, *categoryUUID, db); err != nil {
		response := dto.NewErrorResponse("Category not found", err.Error())
		return c.Status(fiber.StatusNotFound).JSON(response)
	}

	if err := service.DeleteCategory(&category, db); err != nil {
		response := dto.NewErrorResponse("Error deleting category", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	response := dto.NewSuccessResponse(nil, "Category deleted successfully")
	return c.Status(fiber.StatusOK).JSON(response)
}
