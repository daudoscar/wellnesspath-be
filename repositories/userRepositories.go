package repositories

import (
	"wellnesspath/models"

	"gorm.io/gorm"
)

func GetUserByID(tx *gorm.DB, userID uint64) (*models.User, error) {
	var user models.User
	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(tx *gorm.DB, user *models.User) error {
	return tx.Save(user).Error
}

func DeleteUserByID(tx *gorm.DB, userID uint64) error {
	return tx.
		Model(&models.User{}).
		Where("id = ? AND is_deleted = false", userID).
		Update("is_deleted", true).Error
}
