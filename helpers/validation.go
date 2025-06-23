package helpers

import "strings"

var allowedSplits = []string{"Push/Pull/Legs", "Upper/Lower", "Full Body", "Bro Split"}
var allowedGoals = []string{"Muscle Gain", "Fat Loss", "Stamina", "General Fitness"}
var allowedIntensities = []string{"Beginner", "Intermediate", "Advanced"}
var allowedBMICategories = []string{"Underweight", "Normal", "Overweight", "Obese"}
var allowedEquipment = []string{
	"Barbell",
	"Body Only",
	"Cable",
	"Dumbbell",
	"Exercise Ball",
	"Kettlebells",
	"Machine",
	"Medicine Ball",
	"Other",
	"Resistance Bands",
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
