package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/handlers"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func productRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	product := api.Group("/product")
	product.Get("/", handlers.GetProductList)
	product.Get("/:id", handlers.GetProduct)

}
