package service

import (
	"github.com/arsyaputraa/go-synapsis-challenge/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetCategories(products *[]models.Category, query *gorm.DB) error {
	if err := query.Find(products).Error; err != nil {
		return err
	}
	return nil
}

func GetCategoryByID(category *models.Category, categoryID uuid.UUID, db *gorm.DB) error {

	if err := db.First(category, "id = ?", categoryID).Error; err != nil {
		return err
	}
	return nil
}

func CreateCategory(newCategory *models.Category, db *gorm.DB) error {
	if err := db.Create(newCategory).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCategory(category *models.Category, db *gorm.DB) error {
	if err := db.Save(category).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCategory(category *models.Category, db *gorm.DB) error {
	if err := db.Delete(&category).Error; err != nil {
		return err
	}
	return nil
}
