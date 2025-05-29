package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type PlanService struct{}

// GenerateWorkoutPlan creates a personalized workout plan for the user based on their profile.
func (s *PlanService) GenerateWorkoutPlan(userID uint64) (dto.FullPlanOutput, error) {
	// Fetch the user's profile data
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("user profile not found")
	}

	// Prevent duplicate plans: check if the user already has a workout plan
	// existingPlans, err := repositories.GetAllWorkoutPlansByUserID(userID)
	// if err != nil {
	// 	return dto.FullPlanOutput{}, errors.New("failed to check existing workout plans")
	// }
	// if len(existingPlans) > 0 {
	// 	return dto.FullPlanOutput{}, errors.New("user already has a workout plan")
	// }

	// Decode user's equipment preferences from JSON
	equipment := helpers.DecodeEquipment(profile.EquipmentJSON)

	// Get exercises matching user's goal and available equipment
	exercises, err := repositories.GetExercisesByGoalAndEquipment(profile.Goal, equipment)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("failed to fetch matching exercises")
	}
	if len(exercises) == 0 {
		return dto.FullPlanOutput{}, errors.New("no exercises match your profile")
	}

	// Create the base workout plan entry
	plan := models.WorkoutPlan{
		UserID:    userID,
		SplitType: profile.SplitType,
		Goal:      profile.Goal,
	}
	if err := repositories.CreateWorkoutPlan(&plan); err != nil {
		return dto.FullPlanOutput{}, err
	}

	// Get muscle group focus for each day based on the chosen split
	splitFocuses := helpers.GetSplitFocuses(profile.SplitType, profile.Frequency)

	var workoutDays []dto.WorkoutDay // Used for output

	// Generate each day of the plan
	for i, focus := range splitFocuses {
		// Create workout plan day entry
		day := models.WorkoutPlanDay{
			PlanID:    plan.ID,
			DayNumber: i + 1,
			Focus:     focus,
		}
		if err := repositories.CreateWorkoutPlanDay(&day); err != nil {
			return dto.FullPlanOutput{}, err
		}

		// Filter exercises for this day's focus (e.g. chest, legs)
		focused := helpers.FilterExercisesByFocus(exercises, focus)
		if len(focused) == 0 {
			focused = exercises // fallback if no specific match
		}

		// [UPDATE LATER] Randomly select up to 4 exercises
		rand.Shuffle(len(focused), func(i, j int) {
			focused[i], focused[j] = focused[j], focused[i]
		})
		selected := focused
		if len(focused) > 4 {
			selected = focused[:4]
		}

		// Build DTO for the current day
		var dayDTO dto.WorkoutDay
		dayDTO.DayNumber = i + 1
		dayDTO.Focus = focus

		// Create exercises for this day in the plan
		for j, ex := range selected {
			reps := helpers.DetermineReps(profile.Intensity, profile.Goal)

			// Create the plan exercise in the database
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

			// Generate SAS URL for the exercise image
			blobName := "picture/image_" + fmt.Sprint(ex.ID) + ".jpg"
			imageURL, err := helpers.GenerateSASURL(blobName, time.Hour)
			if err != nil {
				return dto.FullPlanOutput{}, fmt.Errorf("failed to generate SAS URL for image %s: %w", blobName, err)
			}

			// Append the exercise to the day's DTO
			dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
				ExerciseID: ex.ID,
				Name:       ex.Name,
				Reps:       reps,
				Sets:       3,
				Order:      j + 1,
				ImageURL:   imageURL,
			})
		}

		workoutDays = append(workoutDays, dayDTO)
	}

	// Final full plan output containing all days + extra info
	output := dto.FullPlanOutput{
		WorkoutPlan:    workoutDays,
		TrainingAdvice: helpers.GenerateTrainingAdvice(profile),
		BMIInfo:        helpers.BuildBMIInfo(profile.BMI, profile.BMICategory),
		CaloriesBurned: helpers.CalculateCalories(profile),
		NutritionPlan:  helpers.GenerateNutrition(profile),
	}

	return output, nil
}

// GetAllPlans returns all workout plans associated with a user (non-deleted)
func (s *PlanService) GetAllPlans(userID uint64) ([]models.WorkoutPlan, error) {
	return repositories.GetAllWorkoutPlansByUserID(userID)
}

// GetPlanByID fetches a specific workout plan and its detailed days & exercises
func (s *PlanService) GetPlanByID(planID uint64) (models.WorkoutPlan, error) {
	return repositories.GetWorkoutPlanWithDetails(planID)
}
