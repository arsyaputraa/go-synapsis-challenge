package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/handlers"
	"github.com/arsyaputraa/go-synapsis-challenge/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
   func userRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api");
	user := api.Group("/user", middleware.JWTMiddleware)
	user.Get("/me", handlers.GetMe)
	user.Patch("/update", handlers.UpdateUser)
	user.Patch("change-password", handlers.UpdatePassword)
   }