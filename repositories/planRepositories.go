package repositories

import (
	"strings"
	"wellnesspath/config"
	"wellnesspath/models"
)

func CreateWorkoutPlan(plan *models.WorkoutPlan) error {
	return config.DB.Create(plan).Error
}

func CreateWorkoutPlanDay(day *models.WorkoutPlanDay) error {
	return config.DB.Create(day).Error
}

func CreateWorkoutPlanExercise(ex *models.WorkoutPlanExercise) error {
	return config.DB.Create(ex).Error
}

// Admin Function (Optional)
func GetAllWorkoutPlansByUserID(userID uint64) ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&plans).Error
	return plans, err
}

func GetWorkoutPlanWithDetails(userID uint64) (models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	err := config.DB.
		Preload("Days.Exercises").
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Order("id").
		First(&plan).Error
	return plan, err
}

func GetExercisesByGoalAndEquipment(goal string, equipmentList []string) ([]models.Exercise, error) {
	var exercises []models.Exercise

	query := config.DB.
		Where("is_deleted = ?", false).
		Where("LOWER(goal_tag) = ? OR LOWER(goal_tag) = ?", strings.ToLower(goal), "general fitness")

	if len(equipmentList) > 0 {
		query = query.Where(buildEquipmentCondition(equipmentList))
	}

	err := query.Find(&exercises).Error
	return exercises, err
}

func DeleteWorkoutPlanByUserID(userID uint64) error {
	return config.DB.
		Model(&models.WorkoutPlan{}).
		Where("user_id = ? AND is_deleted = false", userID).
		Update("is_deleted", true).Error
}

func DeleteFullWorkoutPlanByUserID(userID uint64) error {
	var plans []models.WorkoutPlan
	if err := config.DB.Where("user_id = ? AND is_deleted = 0", userID).Find(&plans).Error; err != nil {
		return err
	}
	for _, plan := range plans {
		if err := config.DB.Where("day_id IN (?)",
			config.DB.Table("workout_plan_days").Select("id").Where("plan_id = ?", plan.ID),
		).Delete(&models.WorkoutPlanExercise{}).Error; err != nil {
			return err
		}
		if err := config.DB.Where("plan_id = ?", plan.ID).Delete(&models.WorkoutPlanDay{}).Error; err != nil {
			return err
		}
		if err := config.DB.Model(&models.WorkoutPlan{}).Where("id = ?", plan.ID).Update("is_deleted", 1).Error; err != nil {
			return err
		}
	}
	return nil
}

func GetActiveWorkoutPlanByUserID(userID uint64) (models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	err := config.DB.
		Preload("Days.Exercises").
		Where("user_id = ? AND is_deleted = ?", userID, false).
		First(&plan).Error
	return plan, err
}

// Helper function to build LIKE OR conditions for equipment matching
func buildEquipmentCondition(equipmentList []string) string {
	var conditions []string
	for _, e := range equipmentList {
		e = strings.ToLower(strings.TrimSpace(e))
		conditions = append(conditions, "LOWER(equipment) LIKE '%"+e+"%'")
	}
	return strings.Join(conditions, " OR ")
}

func FindSimilarExercises(referenceEx models.Exercise, profile *models.Profile, equipment []string, maxCount int) ([]models.Exercise, error) {
	query := config.DB.Model(&models.Exercise{}).
		Where("id != ? AND body_part = ? AND is_deleted = ?", referenceEx.ID, referenceEx.BodyPart, false).
		Where("difficulty <= ?", strings.ToLower(profile.Intensity)).
		Limit(maxCount)

	if len(equipment) > 0 {
		equipCond := buildEquipmentCondition(equipment)
		query = query.Where(equipCond)
	}

	var result []models.Exercise
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func ReplaceExerciseInWorkoutPlan(userID uint64, originalExerciseID uint64, newExerciseID uint64) error {
	var exercise models.WorkoutPlanExercise

	err := config.DB.
		Joins("JOIN workout_plan_days ON workout_plan_days.id = workout_plan_exercises.day_id").
		Joins("JOIN workout_plans ON workout_plans.id = workout_plan_days.plan_id").
		Where("workout_plans.user_id = ? AND workout_plans.is_deleted = false AND workout_plan_exercises.exercise_id = ?", userID, originalExerciseID).
		First(&exercise).Error

	if err != nil {
		return err
	}

	exercise.ExerciseID = newExerciseID
	return config.DB.Save(&exercise).Error
}

func UpdateExerciseInPlanExercise(planExerciseID uint64, newExerciseID uint64) error {
	return config.DB.
		Model(&models.WorkoutPlanExercise{}).
		Where("id = ?", planExerciseID).
		Update("exercise_id", newExerciseID).Error
}
