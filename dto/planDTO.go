package dto

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

type ExercisePlanResponse struct {
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
	DayNumber int                    `json:"dayNumber"`
	Focus     string                 `json:"focus"`
	Exercises []ExercisePlanResponse `json:"exercises"`
}

type FullDayPlanOutput struct {
	WorkoutDay     WorkoutDay `json:"workoutDay"`
	CaloriesBurned float64    `json:"caloriesBurned"`
}
