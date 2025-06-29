package helpers

import (
	"fmt"
	"strings"
	"wellnesspath/dto"
	"wellnesspath/models"
)

func GetSplitFocuses(splitType string, frequency int) []string {
	switch strings.ToLower(splitType) {
	case "push/pull/legs":
		return repeatSequence([]string{"Push", "Pull", "Legs"}, frequency)
	case "upper/lower":
		return repeatSequence([]string{"Upper", "Lower"}, frequency)
	case "full body":
		return repeatSequence([]string{"Full Body"}, frequency)
	case "bro split":
		return repeatSequence([]string{"Chest", "Back", "Legs", "Shoulders", "Arms"}, frequency)
	default:
		return repeatSequence([]string{"General"}, frequency)
	}
}

func repeatSequence(pattern []string, total int) []string {
	var result []string
	for i := 0; i < total; i++ {
		result = append(result, pattern[i%len(pattern)])
	}
	return result
}

func FilterExercisesByFocus(exercises []models.Exercise, focus string) []models.Exercise {
	validParts := GetBodyPartsForFocus(focus)
	var result []models.Exercise

	for _, e := range exercises {
		if Contains(validParts, e.BodyPart) {
			result = append(result, e)
		}
	}

	return result
}

func DetermineReps(intensity, goal, bmiCategory string) int {
	goal = strings.ToLower(goal)
	intensity = strings.ToLower(intensity)
	bmiCategory = strings.ToLower(bmiCategory)

	switch goal {
	case "muscle gain":
		switch intensity {
		case "beginner":
			return 10
		case "intermediate":
			return 8
		default:
			return 6
		}

	case "fat loss":
		if bmiCategory == "overweight" || bmiCategory == "obese" {
			return 15 // more reps, lower load
		}
		return 12

	case "stamina":
		return 15

	case "general fitness":
		if intensity == "beginner" || bmiCategory == "underweight" {
			return 12
		}
		return 10

	default:
		return 10
	}
}

func BuildBMIInfo(bmi float64, category string) dto.BMIInfo {
	advice := map[string]string{
		"Underweight": "Focus on strength and calorie surplus.",
		"Normal":      "Maintain balance across strength and cardio.",
		"Overweight":  "Prioritize fat burning and cardio routines.",
		"Obese":       "Low-impact, high-frequency cardio is recommended.",
	}
	return dto.BMIInfo{
		Value:    bmi,
		Category: category,
		Advice:   advice[category],
	}
}

func GenerateTrainingAdvice(profile *models.Profile) string {
	switch strings.ToLower(profile.Goal) {
	case "muscle gain":
		return "Use moderate to heavy resistance, 8–12 reps per set, with progressive overload. Focus on compound lifts and allow adequate rest between sets (60–90 seconds)."
	case "fat loss":
		return "Use moderate resistance with high volume (12–15 reps), short rest intervals (30–45 seconds), and prioritize supersets or circuits to maximize calorie burn."
	case "stamina":
		return "Train with lighter weights and high reps (15–20+), minimal rest, and maintain steady tempo to build muscular endurance."
	default:
		return "Balance strength and endurance training. Focus on proper form and consistent weekly routines."
	}
}

func CalculateCalories(profile *models.Profile) dto.CaloriesBurned {
	base := float64(profile.DurationPerSession) * 5.0 // average 5 cal/min
	perSession := base
	weekly := perSession * float64(profile.Frequency)
	total := weekly * 4

	return dto.CaloriesBurned{
		PerSession: perSession,
		Weekly:     weekly,
		Total:      total,
	}
}

func getMetValueForGoalTag(goalTag string) float64 {
	// Define MET values based on the GoalTag
	switch goalTag {
	case "Strength":
		return 4.5
	case "Fat Loss":
		return 6.0
	case "General Fitness":
		return 5.0
	case "Stamina":
		return 7.0
	default:
		return 4.0
	}
}

func CalculateTodayCalories(allGoalTags []string, userWeight float64) float64 {
	var totalCalories float64

	// Iterate through all GoalTags
	for _, goalTag := range allGoalTags {
		// Get the MET value for the goalTag
		metValue := getMetValueForGoalTag(goalTag)

		// Calculate calories burned for this GoalTag using a constant time-based factor (e.g., 1 hour)
		// Assuming the exercise lasts for 1 hour
		caloriesBurned := metValue * userWeight // assuming it's per hour

		// Add to total calories
		totalCalories += caloriesBurned
	}

	return totalCalories
}

func GenerateNutrition(profile *models.Profile) dto.DailyNutritionRecommendation {
	// Simple BMR-like logic (Not Fix)
	var activityMultiplier float64
	switch profile.Intensity {
	case "Beginner":
		activityMultiplier = 1.3
	case "Intermediate":
		activityMultiplier = 1.5
	case "Advanced":
		activityMultiplier = 1.7
	default:
		activityMultiplier = 1.4
	}

	calories := int(24 * profile.TargetWeight * activityMultiplier)
	protein := profile.TargetWeight * 1.8

	return dto.DailyNutritionRecommendation{
		Calories: calories,
		Protein:  protein,
	}
}

func GetBodyPartsForFocus(focus string) []string {
	switch strings.ToLower(focus) {
	case "push":
		return []string{"Chest", "Shoulders", "Triceps"}
	case "pull":
		return []string{"Back", "Lats", "Biceps", "Forearms"}
	case "legs":
		return []string{"Quadriceps", "Glutes", "Hamstrings", "Calves", "Lower Back"}
	case "upper":
		return []string{"Chest", "Back", "Shoulders", "Biceps", "Triceps", "Forearms"}
	case "lower":
		return []string{"Quadriceps", "Hamstrings", "Glutes", "Calves", "Lower Back"}
	case "full body":
		return []string{
			"Chest", "Back", "Quadriceps", "Hamstrings", "Calves",
			"Shoulders", "Biceps", "Triceps", "Forearms", "Glutes", "Abdominals",
		}
	case "chest":
		return []string{"Chest"}
	case "back":
		return []string{"Back", "Lats"}
	case "shoulders":
		return []string{"Shoulders"}
	case "arms":
		return []string{"Biceps", "Triceps", "Forearms"}
	default:
		return []string{}
	}
}

func SelectTailoredExercises(exercises []models.Exercise, profile *models.Profile, focus string, maxCount int) []models.Exercise {
	validParts := GetBodyPartsForFocus(focus)
	intensityRank := map[string]int{
		"Beginner":     1,
		"Intermediate": 2,
		"Advanced":     3,
	}
	userRank := intensityRank[strings.ToLower(profile.Intensity)]

	// Tier 1: Strict (match all: bodyPart, goal, difficulty, no dup category)
	selected := filterWithCriteria(exercises, validParts, profile.Goal, userRank, true, true, maxCount)
	if len(selected) >= maxCount {
		return selected
	}

	// Tier 2: Relax goal_tag match (allow all goals)
	selected = filterWithCriteria(exercises, validParts, profile.Goal, userRank, false, true, maxCount)
	if len(selected) >= maxCount {
		return selected
	}

	// Tier 3: Relax difficulty (ignore difficulty level)
	selected = filterWithCriteria(exercises, validParts, profile.Goal, -1, false, false, maxCount)
	if len(selected) >= maxCount {
		return selected
	}

	// Tier 4: Relax category uniqueness
	selected = filterWithCriteria(exercises, validParts, profile.Goal, -1, false, false, maxCount)
	return selected
}

func filterWithCriteria(exs []models.Exercise, validParts []string, goal string, maxRank int, strictGoal bool, uniqueCategory bool, maxCount int) []models.Exercise {
	selected := []models.Exercise{}
	usedCategories := map[string]bool{}

	for _, ex := range exs {
		if !Contains(validParts, ex.BodyPart) {
			continue
		}

		if maxRank != -1 {
			intensityRank := map[string]int{
				"Beginner":     1,
				"Intermediate": 2,
				"Advanced":     3,
			}
			exRank, ok := intensityRank[strings.ToLower(ex.Difficulty)]
			if !ok {
				exRank = 3
			}
			if exRank > maxRank {
				continue
			}
		}

		if strictGoal && !strings.EqualFold(ex.GoalTag, goal) && !strings.EqualFold(ex.GoalTag, "General Fitness") {
			continue
		}

		if uniqueCategory && usedCategories[ex.Category] {
			continue
		}

		selected = append(selected, ex)
		usedCategories[ex.Category] = true

		if len(selected) == maxCount {
			break
		}
	}

	return selected
}

// Helper contains() function
func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

func ValidateSplitAndRestDays(splitType string, frequency int, restDays []int) error {
	if len(restDays) > 7 {
		return fmt.Errorf("invalid number of rest days (max 7)")
	}

	availableDays := 7 - len(restDays)
	if availableDays < frequency {
		return fmt.Errorf("not enough workout days: %d rest days leave only %d workout days, but frequency is set to %d",
			len(restDays), availableDays, frequency)
	}

	switch normalizeSplit(splitType) {
	case "push/pull/legs":
		if frequency < 3 {
			return fmt.Errorf("push/pull/legs split requires at least 3 workout days")
		}
		if frequency%3 != 0 {
			return fmt.Errorf("push/pull/legs split is based on a 3-day cycle; valid frequencies are 3 or 6")
		}
	case "upper/lower":
		if frequency < 2 {
			return fmt.Errorf("upper/lower split requires at least 2 workout days")
		}
		if frequency%2 != 0 {
			return fmt.Errorf("upper/lower split is based on a 2-day cycle; valid frequencies are 2, 4, or 6")
		}
	case "bro split":
		if frequency != 5 {
			return fmt.Errorf("bro split requires exactly 5 workout days")
		}
	case "full body":
		if frequency < 2 {
			return fmt.Errorf("full body split requires at least 2 sessions per week")
		}
		if frequency > 4 {
			return fmt.Errorf("full body split should not exceed 4 sessions per week")
		}
	}

	return nil
}

func GetWorkoutDays(restDays []int, frequency int) ([]int, error) {
	var workoutDays []int
	for i := 1; i <= 7; i++ {
		if !ContainsInt(restDays, i) {
			workoutDays = append(workoutDays, i)
		}
	}
	if len(workoutDays) < frequency {
		return nil, fmt.Errorf("not enough available workout days (got %d, need %d)", len(workoutDays), frequency)
	}
	return workoutDays[:frequency], nil
}

func ContainsInt(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func normalizeSplit(split string) string {
	switch strings.ToLower(split) {
	case "push/pull/legs", "upper/lower", "full body", "bro split":
		return strings.ToLower(split)
	default:
		return "general"
	}
}

func PrepareWorkoutDays(splitType string, frequency int, restDays []int) ([]int, error) {
	if err := ValidateSplitAndRestDays(splitType, frequency, restDays); err != nil {
		return nil, err
	}
	return GetWorkoutDays(restDays, frequency)
}
