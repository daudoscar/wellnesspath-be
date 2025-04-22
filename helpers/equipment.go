package helpers

import (
	"encoding/json"
	"errors"
)

func EncodeEquipment(equipment []string) (string, error) {
	data, err := json.Marshal(equipment)
	if err != nil {
		return "", errors.New("failed to encode equipment")
	}
	return string(data), nil
}

func DecodeEquipment(jsonStr string) []string {
	var result []string
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return []string{}
	}
	return result
}
