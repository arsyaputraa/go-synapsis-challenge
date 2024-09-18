package dto

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
)

type ResponseCategory struct {
	ID          uuid.UUID `json:"id"` // Use UUID as the primary key
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

func NewResponseCategory(c *models.Category) ResponseCategory {

	return ResponseCategory{ID: c.ID, Name: c.Name, Description: c.Description}
}
