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

func GetAllWorkoutPlansByUserID(userID uint64) ([]models.WorkoutPlan, error) {
	var plans []models.WorkoutPlan
	err := config.DB.Where("user_id = ? AND is_deleted = false", userID).Find(&plans).Error
	return plans, err
}

func GetWorkoutPlanWithDetails(planID uint64) (models.WorkoutPlan, error) {
	var plan models.WorkoutPlan
	err := config.DB.
		Preload("Days.Exercises").
		Where("id = ? AND is_deleted = false", planID).
		First(&plan).Error
	return plan, err
}

func GetExercisesByGoalAndEquipment(goal string, equipmentList []string) ([]models.Exercise, error) {
	var exercises []models.Exercise

	// Build base query
	query := config.DB.Where("goal_tag = ? AND is_deleted = false", goal)

	// Apply equipment filtering
	if len(equipmentList) > 0 {
		query = query.Where(buildEquipmentCondition(equipmentList))
	}

	// Run query
	err := query.Find(&exercises).Error
	return exercises, err
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
