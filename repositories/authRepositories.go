package repositories

import (
	"wellnesspath/models"

	"gorm.io/gorm"
)

type AuthRepository struct{}

func InsertUser(tx *gorm.DB, user *models.User) error {
	if err := tx.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(tx *gorm.DB, username string) (models.User, error) {
	var user models.User
	if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
