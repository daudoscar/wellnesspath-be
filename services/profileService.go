package services

import (
	"errors"
	"strings"
	"wellnesspath/config"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type ProfileService struct{}

func (s *ProfileService) GetProfile(userID uint64) (dto.ProfileResponseDTO, error) {
	tx := config.DB.Begin()
	defer tx.Rollback()

	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		return dto.ProfileResponseDTO{}, err
	}

	tx.Commit()

	return dto.ProfileResponseDTO{
		ID:                 profile.ID,
		SplitType:          profile.SplitType,
		Intensity:          profile.Intensity,
		TargetWeight:       profile.TargetWeight,
		BMI:                profile.BMI,
		BMICategory:        profile.BMICategory,
		Frequency:          profile.Frequency,
		DurationPerSession: profile.DurationPerSession,
		Goal:               profile.Goal,
		Equipment:          helpers.DecodeEquipment(profile.EquipmentJSON),
	}, nil
}

func (s *ProfileService) UpdateProfile(userID uint64, input dto.UpdateProfileDTO) error {
	if !helpers.IsValidSplitType(input.SplitType) {
		return errors.New("invalid split type")
	}
	if !helpers.IsValidGoal(input.Goal) {
		return errors.New("invalid goal")
	}
	if !helpers.IsValidIntensity(input.Intensity) {
		return errors.New("invalid intensity")
	}
	if !helpers.IsValidBMICategory(input.BMICategory) {
		return errors.New("invalid BMI category")
	}
	if !helpers.IsValidEquipmentList(input.Equipment) {
		return errors.New("invalid equipment list")
	}

	containsBodyOnly := false
	for _, eq := range input.Equipment {
		if strings.EqualFold(eq, "Body Only") {
			containsBodyOnly = true
			break
		}
	}
	if !containsBodyOnly {
		input.Equipment = append(input.Equipment, "Body Only")
	}

	equipmentJSON, err := helpers.EncodeEquipment(input.Equipment)
	if err != nil {
		return err
	}

	profile := models.Profile{
		UserID:             userID,
		SplitType:          input.SplitType,
		Intensity:          input.Intensity,
		TargetWeight:       input.TargetWeight,
		BMI:                input.BMI,
		BMICategory:        input.BMICategory,
		Frequency:          input.Frequency,
		DurationPerSession: input.DurationPerSession,
		Goal:               input.Goal,
		EquipmentJSON:      equipmentJSON,
	}

	tx := config.DB.Begin()

	existing, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		if err := repositories.CreateProfile(tx, &profile); err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	profile.ID = existing.ID
	if err := repositories.UpdateProfile(tx, &profile); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *ProfileService) DeleteProfile(userID uint64) error {
	tx := config.DB.Begin()

	if err := repositories.DeleteProfileByUserID(tx, userID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
