package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/handlers"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func cartRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	cart := api.Group("/cart", middleware.JWTMiddleware)
	cart.Post("/", handlers.AddToCart)
	cart.Get("/", handlers.GetCartItems)
	cart.Get("/:id", handlers.GetCartItemById)

	cart.Delete("/:id", handlers.RemoveCartItems)
}
