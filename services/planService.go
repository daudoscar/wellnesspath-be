package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type PlanService struct{}

// GenerateWorkoutPlan creates a personalized workout plan for the user based on their profile.
func (s *PlanService) GenerateWorkoutPlan(userID uint64) (dto.FullPlanOutput, error) {
	// Ambil profil user
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("user profile not found")
	}

	// Cek apakah user sudah memiliki workout plan aktif
	existingPlans, err := repositories.GetAllWorkoutPlansByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, fmt.Errorf("failed to check existing workout plans: %w", err)
	}
	if len(existingPlans) > 0 {
		if err := repositories.DeleteFullWorkoutPlanByUserID(userID); err != nil {
			return dto.FullPlanOutput{}, fmt.Errorf("failed to delete existing workout plans: %w", err)
		}
	}

	// Decode equipment user
	equipment := helpers.DecodeEquipment(profile.EquipmentJSON)

	// Ambil semua exercise yang sesuai goal dan equipment
	exercises, err := repositories.GetExercisesByGoalAndEquipment(profile.Goal, equipment)
	if err != nil {
		return dto.FullPlanOutput{}, errors.New("failed to fetch matching exercises")
	}
	if len(exercises) == 0 {
		return dto.FullPlanOutput{}, errors.New("no exercises match your profile")
	}

	// Buat entri utama workout plan
	plan := models.WorkoutPlan{
		UserID:    userID,
		SplitType: profile.SplitType,
		Goal:      profile.Goal,
	}
	if err := repositories.CreateWorkoutPlan(&plan); err != nil {
		return dto.FullPlanOutput{}, err
	}

	// Tentukan fokus setiap hari (Push, Pull, Legs, dst)
	splitFocuses := helpers.GetSplitFocuses(profile.SplitType, profile.Frequency)
	var workoutDays []dto.WorkoutDay

	for i, focus := range splitFocuses {
		// Buat entry WorkoutPlanDay
		day := models.WorkoutPlanDay{
			PlanID:    plan.ID,
			DayNumber: i + 1,
			Focus:     focus,
		}
		if err := repositories.CreateWorkoutPlanDay(&day); err != nil {
			return dto.FullPlanOutput{}, err
		}

		// Filter exercise sesuai fokus hari ini
		focused := helpers.FilterExercisesByFocus(exercises, focus)
		if len(focused) == 0 {
			focused = exercises // fallback jika tidak ada yang match
		}

		// Pilih latihan secara terstruktur
		selected := helpers.SelectTailoredExercises(focused, profile, focus, 4)

		// Fallback 1: Goal diganti "General Fitness"
		if len(selected) == 0 {
			altProfile := *profile
			altProfile.Goal = "General Fitness"
			selected = helpers.SelectTailoredExercises(focused, &altProfile, focus, 4)
		}

		// Fallback 2: Turunkan batas difficulty (contoh: Beginner → Intermediate)
		if len(selected) == 0 && strings.ToLower(profile.Intensity) == "beginner" {
			altProfile := *profile
			altProfile.Intensity = "Intermediate"
			selected = helpers.SelectTailoredExercises(focused, &altProfile, focus, 4)
		}

		// Fallback 3: Ambil berdasarkan body part saja
		if len(selected) == 0 {
			validParts := helpers.GetBodyPartsForFocus(focus)
			selected = []models.Exercise{}
			seen := map[uint64]bool{}
			for _, ex := range focused {
				if helpers.Contains(validParts, ex.BodyPart) && !seen[ex.ID] {
					selected = append(selected, ex)
					seen[ex.ID] = true
					if len(selected) == 4 {
						break
					}
				}
			}
		}

		// Fallback gagal total
		if len(selected) == 0 {
			return dto.FullPlanOutput{}, fmt.Errorf("no suitable exercises found for focus %s", focus)
		}

		// Build response per hari
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

			blobName := "picture/image_" + fmt.Sprint(ex.ID) + ".jpg"
			imageURL, err := helpers.GenerateSASURL(blobName, time.Hour)
			if err != nil {
				return dto.FullPlanOutput{}, fmt.Errorf("failed to generate SAS URL for image %s: %w", blobName, err)
			}

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

	// Build output akhir
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
func (s *PlanService) GetPlanByUserID(userID uint64) (dto.FullPlanOutput, error) {
	// Ambil workout plan aktif milik user
	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, fmt.Errorf("failed to retrieve workout plan: %w", err)
	}

	// Ambil profil user (untuk info BMI, kalori, saran, dll.)
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, fmt.Errorf("failed to retrieve profile: %w", err)
	}

	var workoutDays []dto.WorkoutDay

	for _, day := range plan.Days {
		var dayDTO dto.WorkoutDay
		dayDTO.DayNumber = day.DayNumber
		dayDTO.Focus = day.Focus

		for _, ex := range day.Exercises {
			exerciseDetail, err := repositories.GetExerciseByID(ex.ExerciseID)
			if err != nil {
				continue // skip if missing
			}

			blobName := "picture/image_" + fmt.Sprint(ex.ExerciseID) + ".jpg"
			imageURL, err := helpers.GenerateSASURL(blobName, time.Hour)
			if err != nil {
				imageURL = ""
			}

			dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
				ExerciseID: ex.ExerciseID,
				Name:       exerciseDetail.Name,
				Reps:       ex.Reps,
				Sets:       ex.Sets,
				Order:      ex.Order,
				ImageURL:   imageURL,
			})
		}

		workoutDays = append(workoutDays, dayDTO)
	}

	output := dto.FullPlanOutput{
		WorkoutPlan:    workoutDays,
		TrainingAdvice: helpers.GenerateTrainingAdvice(profile),
		BMIInfo:        helpers.BuildBMIInfo(profile.BMI, profile.BMICategory),
		CaloriesBurned: helpers.CalculateCalories(profile),
		NutritionPlan:  helpers.GenerateNutrition(profile),
	}

	return output, nil
}

func (s *PlanService) DeletePlan(userID uint64) error {
	return repositories.DeleteWorkoutPlanByUserID(userID)
}

func (s *PlanService) GetRecommendedReplacements(userID uint64) ([]dto.ExerciseReplacementResponse, error) {
	profile, err := repositories.GetProfileByUserID(userID)
	if err != nil {
		return nil, err
	}
	equipment := helpers.NormalizeEquipment(profile.EquipmentJSON)

	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return nil, err
	}

	var results []dto.ExerciseReplacementResponse
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			exDetail, err := repositories.GetExerciseByID(ex.ExerciseID)
			if err != nil {
				continue
			}

			similar, err := repositories.FindSimilarExercises(*exDetail, profile, equipment, 3)
			if err != nil {
				continue
			}

			var recs []dto.RecommendedExerciseBrief
			for _, s := range similar {
				recs = append(recs, dto.RecommendedExerciseBrief{
					ExerciseID:  s.ID,
					Name:        s.Name,
					Description: s.Description, // assuming `models.Exercise` has this field
				})
			}

			results = append(results, dto.ExerciseReplacementResponse{
				OriginalExerciseID: ex.ExerciseID,
				Name:               exDetail.Name,
				Replacements:       recs,
			})
		}
	}

	return results, nil
}

// Service function for replacing exercise
func (s *PlanService) ReplaceExercise(userID uint64, req dto.ReplaceExerciseRequest) error {
	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return fmt.Errorf("user has no active workout plan")
	}

	found := false
	var targetPlanExerciseID uint64
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			if ex.ExerciseID == req.OriginalExerciseID {
				targetPlanExerciseID = ex.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return fmt.Errorf("original exercise not found in your plan")
	}

	err = repositories.UpdateExerciseInPlanExercise(targetPlanExerciseID, req.NewExerciseID)
	if err != nil {
		return fmt.Errorf("failed to update exercise: %w", err)
	}

	return nil
}
