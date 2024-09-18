package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/utils"
	"github.com/gofiber/fiber/v2"
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
