package services

import (
	"errors"
	"fmt"
	"time"

	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/repositories"
)

type ExerciseService struct{}

func (s *ExerciseService) GetAllExercises() ([]dto.ExerciseResponseDTO, error) {
	exercises, err := repositories.GetAllExercises()
	if err != nil {
		return nil, err
	}

	var response []dto.ExerciseResponseDTO
	for _, e := range exercises {
		response = append(response, dto.ExerciseResponseDTO{
			ID:                     e.ID,
			Name:                   e.Name,
			BodyPart:               e.BodyPart,
			Difficulty:             e.Difficulty,
			Category:               e.Category,
			ExerciseType:           e.ExerciseType,
			GoalTag:                e.GoalTag,
			Description:            e.Description,
			StepByStepInstructions: e.StepByStepInstructions,
		})
	}

	return response, nil
}

func (s *ExerciseService) GetExerciseByID(id uint64) (dto.ExerciseResponseDTO, error) {
	exercise, err := repositories.GetExerciseByID(id)
	if err != nil {
		return dto.ExerciseResponseDTO{}, errors.New("exercise not found")
	}

	return dto.ExerciseResponseDTO{
		ID:                     exercise.ID,
		Name:                   exercise.Name,
		BodyPart:               exercise.BodyPart,
		Difficulty:             exercise.Difficulty,
		Category:               exercise.Category,
		ExerciseType:           exercise.ExerciseType,
		GoalTag:                exercise.GoalTag,
		Description:            exercise.Description,
		StepByStepInstructions: exercise.StepByStepInstructions,
	}, nil
}

func (s *ExerciseService) GetExerciseVideoByID(id uint64) (dto.VideoResponseDTO, error) {
	blobName := "videos/exercise_" + fmt.Sprint(id) + ".mp4"
	videoURL, err := helpers.GenerateSASURL(blobName, time.Hour)
	if err != nil {
		videoURL = ""
	}

	return dto.VideoResponseDTO{
		ExerciseID: id,
		VideoURL:   videoURL,
	}, nil
}
