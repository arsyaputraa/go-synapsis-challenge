package database

import (
	"log"

	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type SeedData struct {
	Categories []models.Category `json:"categories"`
	Products   []models.Product  `json:"products"`
}

func InitializeAdminUser() {
	var admin models.User
	result := Database.Db.Where("email = ?", "admin@mystore.com").First(&admin)

	if result.Error != nil && result.Error.Error() == "record not found" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("mystore_admin"), 14)
		if err != nil {
			log.Fatalf("Failed to hash password: %v", err)
		}

		admin := models.User{

			Name:     "Admin",
			Email:    "admin@mystore.com",
			Password: string(hashedPassword),
			Role:     models.Admin,
		}

		if err := Database.Db.Create(&admin).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}

		log.Println("Admin user created successfully")
	} else {
		log.Println("Admin user already exists")
	}
}
