package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/handlers"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func adminRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	admin := api.Group("/admin", middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
	admin.Post("/product", handlers.AddProduct)
	admin.Patch("/product/:id", handlers.UpdateProduct)
	admin.Delete("/product/:id", handlers.DeleteProduct)
	admin.Post("/category", handlers.AddCategory)
	admin.Patch("/category/:id", handlers.UpdateCategory)
	admin.Delete("/category/:id", handlers.DeleteCategory)
}
