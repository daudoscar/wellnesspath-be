package helpers

import (
	"strings"
)

func SimilarDifficulty(userIntensity string, exDifficulty string) bool {
	intensityRank := map[string]int{
		"beginner":     1,
		"intermediate": 2,
		"advanced":     3,
	}
	userRank := intensityRank[strings.ToLower(userIntensity)]
	exRank, ok := intensityRank[strings.ToLower(exDifficulty)]
	if !ok {
		exRank = 3 // assume advanced
	}
	return exRank <= userRank
}

func NormalizeEquipment(equipmentJSON string) []string {
	return DecodeEquipment(equipmentJSON)
}
