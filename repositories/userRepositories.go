package repositories

import (
	"wellnesspath/config"
	"wellnesspath/models"
)

func GetUserByID(userID uint64) (*models.User, error) {
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user *models.User) error {
	if err := config.DB.Save(user).Error; err != nil {
		return err
	}
	return nil
}
