package services

import (
	"errors"
	"math/rand"

	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type PlanService struct{}

func (s *PlanService) GenerateWorkoutPlan(userID uint64) (dto.FullPlanOutput, error) {
	// Get user profile
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("user profile not found")
	}
	equipment := helpers.DecodeEquipment(profile.EquipmentJSON)

	// Fetch exercises that match goal + equipment
	exercises, err := repositories.GetExercisesByGoalAndEquipment(profile.Goal, equipment)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("failed to fetch matching exercises")
	}
	if len(exercises) == 0 {
		return dto.FullPlanOutput{}, errors.New("no exercises match your profile")
	}

	// Create workout plan metadata
	plan := models.WorkoutPlan{
		UserID:    userID,
		SplitType: profile.SplitType,
		Goal:      profile.Goal,
	}
	if err := repositories.CreateWorkoutPlan(&plan); err != nil {
		return dto.FullPlanOutput{}, err
	}

	// Create days based on split type
	splitFocuses := helpers.GetSplitFocuses(profile.SplitType, profile.Frequency)
	var workoutDays []dto.WorkoutDay

	for i, focus := range splitFocuses {
		day := models.WorkoutPlanDay{
			PlanID:    plan.ID,
			DayNumber: i + 1,
			Focus:     focus,
		}
		if err := repositories.CreateWorkoutPlanDay(&day); err != nil {
			return dto.FullPlanOutput{}, err
		}

		// Filter exercises by focus
		focused := helpers.FilterExercisesByFocus(exercises, focus)
		rand.Shuffle(len(focused), func(i, j int) {
			focused[i], focused[j] = focused[j], focused[i]
		})
		selected := focused
		if len(focused) > 4 {
			selected = focused[:4]
		}

		var dayDTO dto.WorkoutDay
		dayDTO.DayNumber = i + 1
		dayDTO.Focus = focus

		for j, ex := range selected {
			reps := helpers.DetermineReps(profile.Intensity, profile.Goal)
			planExercise := models.WorkoutPlanExercise{
				DayID:      day.ID,
				ExerciseID: ex.ID,
				Order:      j + 1,
				Reps:       reps,
				Sets:       3,
			}
			if err := repositories.CreateWorkoutPlanExercise(&planExercise); err != nil {
				return dto.FullPlanOutput{}, err
			}

			dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
				ExerciseID: ex.ID,
				Name:       ex.Name,
				Reps:       reps,
				Sets:       3,
				Order:      j + 1,
			})
		}

		workoutDays = append(workoutDays, dayDTO)
	}

	// Build FullPlanOutput DTO
	output := dto.FullPlanOutput{
		WorkoutPlan:    workoutDays,
		TrainingAdvice: helpers.GenerateTrainingAdvice(profile),
		BMIInfo:        helpers.BuildBMIInfo(profile.BMI, profile.BMICategory),
		CaloriesBurned: helpers.CalculateCalories(profile),
		NutritionPlan:  helpers.GenerateNutrition(profile),
	}

	return output, nil
}

func (s *PlanService) GetAllPlans(userID uint64) ([]models.WorkoutPlan, error) {
	return repositories.GetAllWorkoutPlansByUserID(userID)
}

func (s *PlanService) GetPlanByID(planID uint64) (models.WorkoutPlan, error) {
	return repositories.GetWorkoutPlanWithDetails(planID)
}
