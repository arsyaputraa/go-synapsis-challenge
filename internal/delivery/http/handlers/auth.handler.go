package handlers

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/service"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with a name, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param registerDTO body dto.RequestRegister true "User registration data"
// @Success 200 {object} dto.GeneralResponse "Success Message"
// @Failure 400 {object} dto.GeneralResponse "Error Message"
// @Failure 500 {object} dto.GeneralResponse "Error Message"
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	var registerDTO dto.RequestRegister

	// Parse and validate the request body
	if err := c.BodyParser(&registerDTO); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	validate := validator.New()
	if err := validate.Struct(registerDTO); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Use the service layer to register the user
	user, err := service.RegisterUser(registerDTO.Name, registerDTO.Email, registerDTO.Password)
	if err != nil {
		if err == service.ErrEmailExists {
			return c.Status(400).JSON(fiber.Map{"error": "Email already in use"})
		}
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
// @Failure 400 {object} dto.GeneralResponse "Error Message"
// @Failure 404 {object} dto.GeneralResponse "Error Message"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var loginDTO dto.RequestLogin

	// Parse the request body into the struct
	if err := c.BodyParser(&loginDTO); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	validate := validator.New()
	if err := validate.Struct(loginDTO); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Use the service layer to authenticate the user
	user, err := service.AuthenticateUser(loginDTO.Email, loginDTO.Password)
	if err != nil {
		if err == service.ErrUserNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "User not found"})
		} else if err == service.ErrInvalidPassword {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid password"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Could not authenticate user"})
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}
