package repositories

import (
	"fmt"
	"strings"
	"wellnesspath/config"
	"wellnesspath/models"

	"gorm.io/gorm"
)

func CreateWorkoutPlanTx(tx *gorm.DB, plan *models.WorkoutPlan) error {
	return tx.Create(plan).Error
}

func CreateWorkoutPlanDayTx(tx *gorm.DB, day *models.WorkoutPlanDay) error {
	return tx.Create(day).Error
}

func CreateWorkoutPlanExerciseTx(tx *gorm.DB, ex *models.WorkoutPlanExercise) error {
	return tx.Create(ex).Error
}

func CreateWorkoutPlanExercisesBatchTx(tx *gorm.DB, exercises []models.WorkoutPlanExercise) error {
	return tx.Create(&exercises).Error
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

	query = query.Order(`
		CASE 
			WHEN LOWER(equipment) = 'body only' THEN 1 
			ELSE 0 
		END
	`)

	err := query.Find(&exercises).Error
	return exercises, err
}

func DeleteWorkoutPlanByUserID(tx *gorm.DB, userID uint64) error {
	return tx.
		Model(&models.WorkoutPlan{}).
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Update("is_deleted", true).Error
}

func DeleteFullWorkoutPlanByUserIDTx(tx *gorm.DB, userID uint64) error {
	var plans []models.WorkoutPlan
	if err := tx.Where("user_id = ? AND is_deleted = 0", userID).Find(&plans).Error; err != nil {
		return err
	}
	for _, plan := range plans {
		if err := tx.Where("day_id IN (?)",
			tx.Table("workout_plan_days").Select("id").Where("plan_id = ?", plan.ID),
		).Delete(&models.WorkoutPlanExercise{}).Error; err != nil {
			return err
		}
		if err := tx.Where("plan_id = ?", plan.ID).Delete(&models.WorkoutPlanDay{}).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.WorkoutPlan{}).Where("id = ?", plan.ID).Update("is_deleted", 1).Error; err != nil {
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

// HELPER Function
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

func FindExercisesByBodyPartsAndEquipment(
	db *gorm.DB,
	bodyParts []string,
	equipment []string,
	excludeIDs []uint64,
) ([]models.Exercise, error) {
	query := db.Model(&models.Exercise{}).
		Where("is_deleted = ?", false)

	// Filter bodypart
	if len(bodyParts) > 0 {
		query = query.Where("body_part IN ?", bodyParts)
	}

	// Filter equipment
	if len(equipment) > 0 {
		query = query.Where("equipment IN (?)", equipment)
	}

	// Exclude existing exercise IDs
	if len(excludeIDs) > 0 {
		query = query.Where("id NOT IN ?", excludeIDs)
	}

	var result []models.Exercise
	if err := query.Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func UpdateExerciseInPlanExercise(tx *gorm.DB, planExerciseID uint64, newExerciseID uint64) error {
	return tx.
		Model(&models.WorkoutPlanExercise{}).
		Where("id = ?", planExerciseID).
		Update("exercise_id", newExerciseID).Error
}

func UpdateWorkoutPlanExerciseReps(tx *gorm.DB, exerciseID uint64, newReps int) error {
	err := tx.
		Model(&models.WorkoutPlanExercise{}).
		Where("id = ?", exerciseID).
		Update("reps", newReps).Error
	if err != nil {
		return fmt.Errorf("failed to update reps: %w", err)
	}
	return nil
}

func GetWorkoutPlanDayByDayID(planID uint64, dayID uint64) (models.WorkoutPlanDay, error) {
	var day models.WorkoutPlanDay
	err := config.DB.Where("plan_id = ? AND id = ?", planID, dayID).First(&day).Error
	if err != nil {
		return day, fmt.Errorf("workout plan day not found")
	}
	return day, nil
}

func GetWorkoutPlanExerciseByDayID(dayID uint64) ([]models.WorkoutPlanExercise, error) {
	var exercises []models.WorkoutPlanExercise
	err := config.DB.Where("day_id = ?", dayID).Find(&exercises).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercises for day %d", dayID)
	}
	return exercises, nil
}
