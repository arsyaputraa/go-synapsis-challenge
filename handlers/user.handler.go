package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
func GetMe (c *fiber.Ctx) error {
	userID := c.Locals("userID")


    // Fetch the user from the database
    var user models.User
    if err := database.Database.Db.First(&user, "id = ?", userID).Error; err != nil {
		response := dto.NewErrorResponse("User not found", err.Error())
        return c.Status(fiber.StatusNotFound).JSON(response)
    }

    // Return the user details
    response := dto.NewSuccessResponse(user.ToDto(), "User data retrieved successfully")
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
	id := c.Locals("userID")

	var user models.User

    if err := database.Database.Db.First(&user, "id = ?", id).Error; err != nil {
		response := dto.NewErrorResponse("User not found", err.Error())
        return c.Status(fiber.StatusNotFound).JSON(response)
    }


	

	var updateData dto.RequestUpdateUser

	if err := c.BodyParser(&updateData); err != nil {
		response := dto.NewErrorResponse("Invalid Input", err.Error())
        return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	user.Name = updateData.Name
	
	database.Database.Db.Save(&user)

	response := dto.NewSuccessResponse(user.ToDto(), "User Updated!")
    return c.Status(200).JSON(response)

}



// ChangePasswod godoc
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
    // Parse request body
    var updatePasswordDTO dto.RequestUpdatePassword
    if err := c.BodyParser(&updatePasswordDTO); err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("Cannot parse JSON", err.Error()))
    }

    // Validate the input
    if err := validate.Struct(updatePasswordDTO); err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("Validation error", err.Error()))
    }

    // Get the current user's ID from the context
    userID := c.Locals("userID").(uuid.UUID)

    // Find the user in the database
    var user models.User
    if err := database.Database.Db.First(&user, "id = ?", userID).Error; err != nil {
        return c.Status(404).JSON(dto.NewErrorResponse("User not found", err.Error()))
    }

    // Verify the current password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updatePasswordDTO.CurrentPassword)); err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("Incorrect current password", nil))
    }

    // Hash the new password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatePasswordDTO.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(500).JSON(dto.NewErrorResponse("Failed to hash the new password", err.Error()))
    }

    // Update the user's password in the database
    user.Password = string(hashedPassword)
    if err := database.Database.Db.Save(&user).Error; err != nil {
        return c.Status(500).JSON(dto.NewErrorResponse("Failed to update password", err.Error()))
    }

    return c.JSON(dto.NewSuccessResponse(nil, "Password updated successfully"))
}



