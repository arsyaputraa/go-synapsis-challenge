package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetMe godoc
// @Summary Get current authenticated user
// @Description Get the details of the currently authenticated user
// @Tags user
// @Produce  json
// @Success 200 {object} dto.GeneralResponse "User data retrieved successfully"
// @Failure 404 {object} dto.GeneralResponse "User not found"
// @Failure 401 {object} dto.GeneralResponse "Unauthorized"
// @Router /user/me [get]
// @Security BearerAuth
func GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	// Fetch the user using the service
	user, err := service.GetUserByID(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("User not found", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error retrieving user", err.Error()))
	}

	// Return the user details
	response := dto.NewSuccessResponse(dto.NewResponseUser(user), "User data retrieved successfully")
	return c.JSON(response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user detail
// @Tags user
// @Accept  json
// @Produce  json
// @Param UpdateData body dto.RequestUpdateUser true "User update data"
// @Success 200 {object} dto.GeneralResponse "Success Message"
// @Failure 404 {object} dto.GeneralResponse "Error Message"
// @Failure 401 {object} dto.GeneralResponse "Error Message"
// @Router /user/update [patch]
// @Security BearerAuth
func UpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	var updateData dto.RequestUpdateUser
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse("Invalid Input", err.Error()))
	}

	// Use the service layer to update the user
	user, err := service.UpdateUserDetails(userID, updateData.Name)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse("User not found", err.Error()))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse("Error updating user", err.Error()))
	}

	response := dto.NewSuccessResponse(dto.NewResponseUser(user), "User Updated!")
	return c.Status(200).JSON(response)
}

// ChangePassword godoc
// @Summary Update Password
// @Description Update user Password
// @Tags user
// @Accept  json
// @Produce  json
// @Param NewPassword body dto.RequestUpdatePassword true "User update Password Data"
// @Success 200 {object} dto.GeneralResponse "Success Message"
// @Failure 404 {object} dto.GeneralResponse "Error Message"
// @Failure 401 {object} dto.GeneralResponse "Error Message"
// @Router /user/change-password [patch]
// @Security BearerAuth
func UpdatePassword(c *fiber.Ctx) error {
	var updatePasswordDTO dto.RequestUpdatePassword
	if err := c.BodyParser(&updatePasswordDTO); err != nil {
		return c.Status(400).JSON(dto.NewErrorResponse("Cannot parse JSON", err.Error()))
	}

	validate := validator.New()
	if err := validate.Struct(updatePasswordDTO); err != nil {
		return c.Status(400).JSON(dto.NewErrorResponse("Validation error", err.Error()))
	}

	// Get the current user's ID from the context
	userID := c.Locals("userID").(uuid.UUID)

	// Use the service layer to update the password
	err := service.UpdateUserPassword(userID, updatePasswordDTO.CurrentPassword, updatePasswordDTO.NewPassword)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(404).JSON(dto.NewErrorResponse("User not found", err.Error()))
		} else if err == service.ErrInvalidCurrentPassword {
			return c.Status(400).JSON(dto.NewErrorResponse("Incorrect current password", nil))
		}
		return c.Status(500).JSON(dto.NewErrorResponse("Failed to update password", err.Error()))
	}

	return c.JSON(dto.NewSuccessResponse(nil, "Password updated successfully"))
}
