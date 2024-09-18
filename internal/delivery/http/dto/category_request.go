package dto

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
)

type RequestCategory struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type RequestUpdateCategory struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (rp *RequestCategory) ToModel() models.Category {
	return models.Category{
		Name:        rp.Name,
		Description: rp.Description,
	}
}
