package router

import (
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
   func SetupRoutes(app *fiber.App) {
	// user
	userRoutes(app)
	authRoutes(app)

   }

