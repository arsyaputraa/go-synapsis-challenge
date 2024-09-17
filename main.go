package main

import (
	"log"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	_ "github.com/arsyaputraa/go-synapsis-challenge/docs"
	"github.com/arsyaputraa/go-synapsis-challenge/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Online Store API
// @version 1.0
// @description This is an online store API using Golang, Fiber, and GORM.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "Bearer <token>"
func init() {

	database.Connect();
}


func main() {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/swagger/*", swagger.HandlerDefault)

	router.SetupRoutes(app)

	log.Fatal(app.Listen(":8080"))
}