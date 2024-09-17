package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/arsyaputraa/go-synapsis-challenge/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)



var validate = validator.New()


// Register godoc
// @Summary Register a new user
// @Description Register a new user with a name, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param registerDTO body dto.RequestRegister true "User registration data"
/// @Success 200 {object} dto.GeneralResponse "Success Message"
// @Failure 404 {object} dto.GeneralResponse "Error Message"
// @Failure 401 {object} dto.GeneralResponse "Error Message"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
    var registerDTO dto.RequestRegister

    // Parse and validate the request body
    if err := c.BodyParser(&registerDTO); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Validate the struct
    if err := validate.Struct(registerDTO); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    // Hash the password
    passwordHash, _ := bcrypt.GenerateFromPassword([]byte(registerDTO.Password), 14)

    // Create the user with a UUID
    user := models.User{
        Name: registerDTO.Name,
        Email:    registerDTO.Email,
        Password: string(passwordHash),
        Role: models.Customer,
    }

    // Save the user in the database
    if err := database.Database.Db.Create(&user).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Could not register user"})
    }

    return c.Status(201).JSON(fiber.Map{"message": "User registered successfully", "id": user.ID})
}




// Login godoc
// @Summary Log in a user
// @Description Authenticate a user with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param loginDTO body dto.RequestLogin true "User login data"
// @Success 200 {object} dto.GeneralResponse "Success Message"
// @Failure 404 {object} dto.GeneralResponse "Error Message"
// @Failure 401 {object} dto.GeneralResponse "Error Message"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
    var loginDTO dto.RequestLogin

    // Parse the request body into the struct
    if err := c.BodyParser(&loginDTO); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Validate the struct
    if err := validate.Struct(loginDTO); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": err.Error()})
    }

    // Find the user by email
    var user models.User
    if err := database.Database.Db.Where("email = ?", loginDTO.Email).First(&user).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "User not found"})
    }

    // Compare the hashed password with the provided password
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDTO.Password)); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid password"})
    }

    // Generate JWT token
    token, err := utils.GenerateJWT(user.ID, user.Role)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
    }

    return c.JSON(fiber.Map{"token": token})
}


