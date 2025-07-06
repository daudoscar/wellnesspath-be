package dto

import "wellnesspath/models"

type FullPlanOutput struct {
	WorkoutPlan    []WorkoutDay                 `json:"workoutPlan"`
	Schedule       []ScheduledExercise          `json:"schedule,omitempty"`
	BMIInfo        BMIInfo                      `json:"bmiInfo"`
	CaloriesBurned CaloriesBurned               `json:"caloriesBurned"`
	NutritionPlan  DailyNutritionRecommendation `json:"nutritionPlan"`
	TrainingAdvice string                       `json:"trainingAdvice"`
}

type WorkoutDay struct {
	DayNumber int                    `json:"dayNumber"`
	Focus     string                 `json:"focus"`
	Exercises []ExercisePlanResponse `json:"exercises"`
}

type WorkoutDayToday struct {
	DayNumber int                     `json:"dayNumber"`
	Focus     string                  `json:"focus"`
	Exercises []ExerciseTodayResponse `json:"exercises"`
}

type ExercisePlanResponse struct {
	ExerciseID uint64 `json:"exerciseId"`
	Name       string `json:"name"`
	Reps       int    `json:"reps"`
	Sets       int    `json:"sets"`
	Order      int    `json:"order"`
	Note       string `json:"note,omitempty"`
	BodyPart   string `json:"body_part"`
	Equipment  string `json:"equipment"`
}

type ExerciseTodayResponse struct {
	ExerciseID uint64 `json:"exerciseId"`
	Name       string `json:"name"`
	Reps       int    `json:"reps"`
	Sets       int    `json:"sets"`
	Order      int    `json:"order"`
	Note       string `json:"note,omitempty"`
	ImageURL   string `json:"image_url"`
}

type ScheduledExercise struct {
	DayNumber int    `json:"dayNumber"`
	Date      string `json:"date"`
	Exercise  string `json:"exercise"`
	Reps      int    `json:"reps"`
	Sets      int    `json:"sets"`
}

type BMIInfo struct {
	Value    float64 `json:"value"`
	Category string  `json:"category"`
	Advice   string  `json:"advice"`
}

type CaloriesBurned struct {
	PerSession float64 `json:"perSession"`
	Weekly     float64 `json:"weekly"`
	Total      float64 `json:"total"`
}

type NutritionPlan struct {
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fats     float64 `json:"fats"`
}

type ReplacementExercise struct {
	ExerciseID uint64 `json:"exerciseId"`
	Name       string `json:"name"`
	ImageURL   string `json:"image_url"`
}

type ExerciseReplacementResponse struct {
	OriginalExerciseID uint64                     `json:"originalExerciseId"`
	Name               string                     `json:"name"`
	Replacements       []RecommendedExerciseBrief `json:"replacements"`
}

type RecommendedExerciseBrief struct {
	ExerciseID  uint64 `json:"exerciseId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ReplaceExerciseRequest struct {
	OriginalExerciseID uint64 `json:"originalExerciseID" binding:"required"`
	NewExerciseID      uint64 `json:"newExerciseID" binding:"required"`
}

type EditRepsRequest struct {
	PlanExerciseID uint64 `json:"planExerciseId" binding:"required"`
	NewReps        int    `json:"newReps" binding:"required,min=1,max=100"`
}

type WorkoutDayOutput struct {
	DayNumber int                     `json:"dayNumber"`
	Focus     string                  `json:"focus"`
	Exercises []ExerciseTodayResponse `json:"exercises"`
}

type FullDayPlanOutput struct {
	WorkoutDay     WorkoutDayToday `json:"workoutDay"`
	CaloriesBurned float64         `json:"caloriesBurned"`
}

type InitializePlanOutput struct {
	PlanID    uint64            `json:"planId"`
	Profile   models.Profile    `json:"profile"`
	RestDays  []int             `json:"restDays"`
	Exercises []models.Exercise `json:"exercises"`
}

type CreateDaysRequest struct {
	PlanID   uint64         `json:"plan_id"`
	RestDays []int          `json:"rest_days"`
	Profile  models.Profile `json:"profile"`
}

type InsertExercisesRequest struct {
	Profile models.Profile          `json:"profile"`
	Days    []models.WorkoutPlanDay `json:"days"`
}
