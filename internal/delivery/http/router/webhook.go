package router

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/handlers"
	"github.com/gofiber/fiber/v2"
)

func webhookRoutes(app *fiber.App) {
	// grouping
	api := app.Group("/api")
	webhook := api.Group("/webhook")
	webhook.Get("/payment", handlers.PaymentWebhook)

}
