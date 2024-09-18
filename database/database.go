package database

import (
	"fmt"
	"log"
	"os"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/arsyaputraa/go-synapsis-challenge/pkg/config"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database instance
type Dbinstance struct {
	Db *gorm.DB
}

var Database Dbinstance

// Connect function
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.Config("DB_HOST"),
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
		config.Config("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}
	log.Println("Connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Cart{}, &models.CartItem{}, &models.Order{}, &models.OrderItem{}, &models.Payment{})
	Database = Dbinstance{
		Db: db,
	}
}
