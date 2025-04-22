package helpers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var allowedSplits = []string{"Push/Pull/Legs", "Upper/Lower", "Full Body", "Bro Split"}
var allowedGoals = []string{"Muscle Gain", "Fat Loss", "Stamina", "General Fit"}
var allowedIntensities = []string{"Beginner", "Intermediate", "Advanced"}
var allowedBMICategories = []string{"Underweight", "Normal", "Overweight", "Obese"}

var allowedEquipment = []string{
	"Body Only",
	"Bands",
	"Barbell",
	"Cable",
	"Cables",
	"Dumbbell",
	"Dumbbells",
	"E-Z Curl Bar",
	"Exercise Ball",
	"Kettlebells",
	"Machine",
	"Medicine Ball",
	"Weight Bench",
	"None",
	"Other",
}

func SplitTypeValidator(fl validator.FieldLevel) bool {
	return containsCaseInsensitive(allowedSplits, fl.Field().String())
}

func GoalValidator(fl validator.FieldLevel) bool {
	return containsCaseInsensitive(allowedGoals, fl.Field().String())
}

func IntensityValidator(fl validator.FieldLevel) bool {
	return containsCaseInsensitive(allowedIntensities, fl.Field().String())
}

func BMICategoryValidator(fl validator.FieldLevel) bool {
	return containsCaseInsensitive(allowedBMICategories, fl.Field().String())
}

func EquipmentValidator(fl validator.FieldLevel) bool {
	equipmentList, ok := fl.Field().Interface().([]string)
	if !ok {
		return false
	}
	for _, item := range equipmentList {
		if !containsCaseInsensitive(allowedEquipment, item) {
			return false
		}
	}
	return true
}

func containsCaseInsensitive(list []string, input string) bool {
	for _, v := range list {
		if strings.EqualFold(v, input) {
			return true
		}
	}
	return false
}
