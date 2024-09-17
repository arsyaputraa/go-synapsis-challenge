package middleware

import (
	"log"
	"strings"

	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/arsyaputraa/go-synapsis-challenge/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// JWTMiddleware validates the JWT and extracts the user ID
func JWTMiddleware(c *fiber.Ctx) error {
    // Get the token from the Authorization header
	authHeader := c.Get("Authorization")

    if authHeader == "" {
        response := dto.NewErrorResponse("Unauthorized", nil)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }

    // Check if the header contains 'Bearer' and the token
    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
        response := dto.NewErrorResponse("Invalid JWT format, make sure to include Bearer", nil)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }
    tokenString := parts[1]

 

    // Parse the JWT token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return utils.GetJWTSecret(), nil
    })
    if err != nil || !token.Valid {
        response := dto.NewErrorResponse("Unauthorized", nil)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }

    // Extract user ID from the token
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        response := dto.NewErrorResponse("Unauthorized", nil)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }
	

    userID, err := uuid.Parse(claims["user_id"].(string))
    if err != nil {
        response := dto.NewErrorResponse("Unauthorized", err)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }

    role, ok := claims["role"].(string)
    if !ok {
        response := dto.NewErrorResponse("Unauthorized", err)
        return c.Status(fiber.StatusUnauthorized).JSON(response)
    }

	log.Printf("user role: %s", role)
    // Store the user ID in the context
    c.Locals("userID", userID)
    c.Locals("role",role )

    // Proceed to the next handler
    return c.Next()
}