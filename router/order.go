package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/handlers"
	"github.com/arsyaputraa/go-synapsis-challenge/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func orderRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	order := api.Group("/order", middleware.JWTMiddleware)

	order.Post("/checkout", handlers.CheckoutOrder)

}
