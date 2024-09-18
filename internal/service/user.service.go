package service

import (
	"errors"

	"github.com/arsyaputraa/go-synapsis-challenge/database"
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCurrentPassword = errors.New("incorrect current password")
)

// GetUserByID retrieves a user by their ID.
func GetUserByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := database.Database.Db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

// UpdateUserDetails updates user details like name.
func UpdateUserDetails(userID uuid.UUID, name string) (*models.User, error) {
	var user models.User
	if err := database.Database.Db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, ErrUserNotFound
	}

	user.Name = name

	if err := database.Database.Db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUserPassword updates the password for the user.
func UpdateUserPassword(userID uuid.UUID, currentPassword, newPassword string) error {
	var user models.User

	// Find the user by ID
	if err := database.Database.Db.First(&user, "id = ?", userID).Error; err != nil {
		return ErrUserNotFound
	}

	// Verify the current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return ErrInvalidCurrentPassword
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password
	user.Password = string(hashedPassword)
	if err := database.Database.Db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
