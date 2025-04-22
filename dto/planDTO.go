package dto

type FullPlanOutput struct {
	WorkoutPlan    []WorkoutDay        `json:"workoutPlan"`
	Schedule       []ScheduledExercise `json:"schedule,omitempty"` // optional future use
	BMIInfo        BMIInfo             `json:"bmiInfo"`
	CaloriesBurned CaloriesBurned      `json:"caloriesBurned"`
	NutritionPlan  NutritionPlan       `json:"nutritionPlan"`
	TrainingAdvice string              `json:"trainingAdvice"`
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
}

type ScheduledExercise struct {
	DayNumber int    `json:"dayNumber"`
	Date      string `json:"date"` // e.g. "2025-04-26"
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
