package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/handlers"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
   func authRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api");
	auth := api.Group("/auth")
	auth.Post("/register", handlers.Register)
	auth.Post("/login", handlers.Login)
   }