package service

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"gorm.io/gorm"
)

func GetCategories(products *[]models.Category, query *gorm.DB) error {
	if err := query.Find(products).Error; err != nil {
		return err
	}
	return nil
}
