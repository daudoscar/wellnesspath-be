package services

import (
	"errors"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type ProfileService struct{}

func (s *ProfileService) GetProfile(userID uint64) (dto.ProfileResponseDTO, error) {
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return dto.ProfileResponseDTO{}, err
	}

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

	existing, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return repositories.CreateProfile(&profile)
	}

	profile.ID = existing.ID
	return repositories.UpdateProfile(&profile)
}
