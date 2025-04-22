package repositories

import (
	"syncspend/config"
	"syncspend/models"
)

type AuthRepository struct{}

func InsertUser(user *models.User) error {
	if err := config.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}
