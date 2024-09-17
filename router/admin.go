package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/handlers"
	"github.com/arsyaputraa/go-synapsis-challenge/middleware"
	"github.com/gofiber/fiber/v2"
)

func adminRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	admin := api.Group("/admin", middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	admin.Post("/product", handlers.AddProduct)
	admin.Patch("/product/:id", handlers.UpdateProduct)
	admin.Delete("/product/:id", handlers.DeleteProduct)
}
