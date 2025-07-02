package dto

type ExerciseResponseDTO struct {
	ID                     uint64 `json:"id"`
	Name                   string `json:"name"`
	BodyPart               string `json:"body_part"`
	Difficulty             string `json:"difficulty"`
	Category               string `json:"category"`
	ExerciseType           string `json:"exercise_type"`
	GoalTag                string `json:"goal_tag"`
	Description            string `json:"description"`
	StepByStepInstructions string `json:"step_by_step_instructions"`
	Equipment              string `json:"equipment"`
}

type VideoResponseDTO struct {
	ExerciseID uint64 `json:"exercise_id"`
	VideoURL   string `json:"video_url"`
}
