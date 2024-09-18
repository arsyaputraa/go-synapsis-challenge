package service

import (
	"errors"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("wrong email or password")
	ErrInvalidPassword = errors.New("wrong email or password")
	ErrEmailExists     = errors.New("email already exists")
)

// RegisterUser registers a new user in the database.
func RegisterUser(name, email, password string) (*models.User, error) {
	// Check if the user already exists
	var existingUser models.User
	if err := database.Database.Db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, ErrEmailExists
	}

	// Hash the password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, err
	}

	// Create the user
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: string(passwordHash),
		Role:     models.Customer,
	}

	// Save the user in the database
	if err := database.Database.Db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user with email and password.
func AuthenticateUser(email, password string) (*models.User, error) {
	var user models.User

	// Find the user by email
	if err := database.Database.Db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// Compare the hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidPassword
	}

	return &user, nil
}
