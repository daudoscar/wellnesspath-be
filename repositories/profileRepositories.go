package repositories

import (
	"errors"
	"wellnesspath/models"

	"gorm.io/gorm"
)

func GetProfileByUserID(tx *gorm.DB, userID uint64) (*models.Profile, error) {
	var profile models.Profile
	if err := tx.Where("user_id = ? AND is_deleted = ?", userID, false).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func CreateProfile(tx *gorm.DB, profile *models.Profile) error {
	return tx.Create(profile).Error
}

func UpdateProfile(tx *gorm.DB, profile *models.Profile) error {
	return tx.Save(profile).Error
}

func DeleteProfileByUserID(tx *gorm.DB, userID uint64) error {
	return tx.
		Model(&models.Profile{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Update("is_deleted", true).Error
}
