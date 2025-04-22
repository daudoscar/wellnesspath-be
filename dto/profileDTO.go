package dto

type UpdateProfileDTO struct {
	SplitType          string   `json:"split_type" binding:"required,split"`
	Intensity          string   `json:"intensity" binding:"required,intensity"`
	TargetWeight       float64  `json:"target_weight"`
	BMI                float64  `json:"bmi"`
	BMICategory        string   `json:"bmi_category" binding:"required,bmi_category"`
	Frequency          int      `json:"frequency"`
	DurationPerSession int      `json:"duration_per_session"`
	Goal               string   `json:"goal" binding:"required,goal"`
	Equipment          []string `json:"equipment" binding:"equipment"`
}

type ProfileResponseDTO struct {
	ID                 uint64   `json:"id"`
	SplitType          string   `json:"split_type"`
	Intensity          string   `json:"intensity"`
	TargetWeight       float64  `json:"target_weight"`
	BMI                float64  `json:"bmi"`
	BMICategory        string   `json:"bmi_category"`
	Frequency          int      `json:"frequency"`
	DurationPerSession int      `json:"duration_per_session"`
	Goal               string   `json:"goal"`
	Equipment          []string `json:"equipment"`
}
