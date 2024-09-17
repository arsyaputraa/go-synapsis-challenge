package dto

import (
	"github.com/arsyaputraa/go-synapsis-challenge/models"
	"github.com/google/uuid"
)

type ResponseCategory struct {
	ID          uuid.UUID `json:"id,omitempty"` // Use UUID as the primary key
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
}

func NewResponseCategory(c *models.Category) ResponseCategory {

	return ResponseCategory{ID: c.ID, Name: c.Name, Description: c.Description}
}
