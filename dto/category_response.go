package dto

import "github.com/google/uuid"

type ResponseCategory struct {
	ID          uuid.UUID `json:"id,omitempty"` // Use UUID as the primary key
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}