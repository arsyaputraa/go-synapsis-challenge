package middleware

import (
	"github.com/arsyaputraa/go-synapsis-challenge/dto"
	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(requiredRole string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string) // Extract role from context (set in JWT middleware)

		if role != requiredRole {
		response := dto.NewErrorResponse("Forbidden",nil)
        return c.Status(fiber.StatusForbidden).JSON(response)
		}

		return c.Next()
	}
}