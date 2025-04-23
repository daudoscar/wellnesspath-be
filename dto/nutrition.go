package dto

type DailyNutritionRecommendation struct {
	Calories int     `json:"calories"`
	Protein  float64 `json:"protein"` // in grams
}
