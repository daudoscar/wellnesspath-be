package repositories

import (
	"wellnesspath/config"
	"wellnesspath/models"
)

func GetAllExercises() ([]models.Exercise, error) {
	var exercises []models.Exercise
	if err := config.DB.Where("is_deleted = ?", false).Find(&exercises).Error; err != nil {
		return nil, err
	}
	return exercises, nil
}

func GetExerciseByID(id uint64) (*models.Exercise, error) {
	var exercise models.Exercise
	if err := config.DB.Where("id = ? AND is_deleted = ?", id, false).First(&exercise).Error; err != nil {
		return nil, err
	}
	return &exercise, nil
}

func GetExercisesByIDs(ids []uint64) (map[uint64]*models.Exercise, error) {
	var exercises []models.Exercise
	if err := config.DB.Where("id IN ?", ids).Find(&exercises).Error; err != nil {
		return nil, err
	}

	exerciseMap := make(map[uint64]*models.Exercise)
	for i := range exercises {
		exercise := exercises[i]
		exerciseMap[exercise.ID] = &exercise
	}

	return exerciseMap, nil
}
