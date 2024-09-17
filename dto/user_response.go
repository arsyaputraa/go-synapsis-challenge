package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
)

type ResponseUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      string    `json:"role"`
}

func NewResponseUser(u *models.User) ResponseUser {
	return ResponseUser{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, Role: u.Role}
}
