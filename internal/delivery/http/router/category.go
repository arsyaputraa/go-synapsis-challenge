package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func categoryRoutes(app *fiber.App) {
	api := app.Group("/api")
	category := api.Group("/category")
	category.Get("/", handlers.GetCategories)

}
