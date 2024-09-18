package main

import (
	"log"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	_ "github.com/arsyaputraa/go-synapsis-challenge/docs"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/delivery/http/router"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Online Store API
// @version 1.0
// @description This is an online store API using Golang, Fiber, and GORM.

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Bearer <token>"
func init() {
	config.LoadEnv()
	database.Connect()
}

func main() {
	database.InitializeAdminUser()

	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/swagger/*", swagger.HandlerDefault)
	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8080"))
}
