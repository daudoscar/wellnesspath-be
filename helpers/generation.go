package helpers

import (
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
	var result []models.Exercise
	for _, e := range exercises {
		if strings.EqualFold(e.BodyPart, focus) || strings.Contains(strings.ToLower(e.Category), strings.ToLower(focus)) {
			result = append(result, e)
		}
	}
	return result
}

func DetermineReps(intensity, goal string) int {
	switch strings.ToLower(goal) {
	case "muscle gain":
		if intensity == "Beginner" {
			return 10
		} else if intensity == "Intermediate" {
			return 8
		} else {
			return 6
		}
	case "fat loss":
		return 12
	case "stamina":
		return 15
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
		return "Focus on progressive overload. Increase weights weekly."
	case "fat loss":
		return "Include more HIIT or circuits. Keep rest short."
	case "stamina":
		return "Emphasize continuous movement. Minimize rest."
	default:
		return "Stay consistent and listen to your body."
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

func GenerateNutrition(profile *models.Profile) dto.NutritionPlan {
	calories := 2000.0
	if strings.ToLower(profile.Goal) == "muscle gain" {
		calories += 300
	} else if strings.ToLower(profile.Goal) == "fat loss" {
		calories -= 300
	}

	return dto.NutritionPlan{
		Calories: calories,
		Protein:  calories * 0.3 / 4,
		Carbs:    calories * 0.4 / 4,
		Fats:     calories * 0.3 / 9,
	}
}
