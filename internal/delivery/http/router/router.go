package router

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {
	// user
	userRoutes(app)
	// auth
	authRoutes(app)
	// product
	productRoutes(app)
	// admin
	adminRoutes(app)
	// cart
	cartRoutes(app)
	// order
	orderRoutes(app)
	// category
	categoryRoutes(app)
	// webhook
	webhookRoutes(app)
}
