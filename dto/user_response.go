package dto

import (
	"time"

	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
)

type ResponseUser struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      string    `json:"role,omitempty"`
}

func NewResponseUser(u *models.User) ResponseUser {
	return ResponseUser{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt, Role: u.Role}
}
