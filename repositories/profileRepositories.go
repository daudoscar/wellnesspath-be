package repositories

import (
	"errors"
	"wellnesspath/config"
	"wellnesspath/models"

	"gorm.io/gorm"
)

func GetProfileByUserID(userID uint64) (*models.Profile, error) {
	var profile models.Profile
	if err := config.DB.Where("user_id = ? AND is_deleted = false", userID).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &profile, nil
}

func CreateProfile(profile *models.Profile) error {
	return config.DB.Create(profile).Error
}

func UpdateProfile(profile *models.Profile) error {
	return config.DB.Save(profile).Error
}
