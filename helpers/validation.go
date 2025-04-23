package helpers

import "strings"

var allowedSplits = []string{"Push/Pull/Legs", "Upper/Lower", "Full Body", "Bro Split"}
var allowedGoals = []string{"Muscle Gain", "Fat Loss", "Stamina", "General Fit"}
var allowedIntensities = []string{"Beginner", "Intermediate", "Advanced"}
var allowedBMICategories = []string{"Underweight", "Normal", "Overweight", "Obese"}
var allowedEquipment = []string{
	"Body Only", "Bands", "Barbell", "Cable", "Cables", "Dumbbell", "Dumbbells",
	"E-Z Curl Bar", "Exercise Ball", "Kettlebells", "Machine", "Medicine Ball",
	"Weight Bench", "None", "Other",
}

// Split
func IsValidSplitType(value string) bool {
	return containsCaseInsensitive(allowedSplits, value)
}

// Goal
func IsValidGoal(value string) bool {
	return containsCaseInsensitive(allowedGoals, value)
}

// Intensity
func IsValidIntensity(value string) bool {
	return containsCaseInsensitive(allowedIntensities, value)
}

// BMI Category
func IsValidBMICategory(value string) bool {
	return containsCaseInsensitive(allowedBMICategories, value)
}

// Equipment list
func IsValidEquipmentList(equipmentList []string) bool {
	for _, eq := range equipmentList {
		if !containsCaseInsensitive(allowedEquipment, eq) {
			return false
		}
	}
	return true
}

// Helper
func containsCaseInsensitive(list []string, input string) bool {
	for _, v := range list {
		if strings.EqualFold(v, input) {
			return true
		}
	}
	return false
}
