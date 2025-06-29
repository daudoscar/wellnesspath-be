package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"wellnesspath/config"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/models"
	"wellnesspath/repositories"
)

type PlanService struct{}

// GenerateWorkoutPlan creates a personalized workout plan for the user based on their profile.
func (s *PlanService) GenerateWorkoutPlan(userID uint64) error {
	tx := config.DB.Begin()

	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return errors.New("user profile not found")
	}

	var restDays []int
	if err := json.Unmarshal([]byte(profile.RestDaysJSON), &restDays); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to parse rest days: %w", err)
	}

	if err := helpers.ValidateSplitAndRestDays(profile.SplitType, profile.Frequency, restDays); err != nil {
		tx.Rollback()
		return err
	}

	existingPlans, err := repositories.GetAllWorkoutPlansByUserID(userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check existing workout plans: %w", err)
	}
	if len(existingPlans) > 0 {
		if err := repositories.DeleteFullWorkoutPlanByUserIDTx(tx, userID); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to delete existing workout plans: %w", err)
		}
	}

	equipment := helpers.DecodeEquipment(profile.EquipmentJSON)
	exercises, err := repositories.GetExercisesByGoalAndEquipment(profile.Goal, equipment)
	if err != nil {
		tx.Rollback()
		return errors.New("failed to fetch matching exercises")
	}
	if len(exercises) == 0 {
		tx.Rollback()
		return errors.New("no exercises match your profile")
	}

	plan := models.WorkoutPlan{
		UserID:    userID,
		SplitType: profile.SplitType,
		Goal:      profile.Goal,
	}
	if err := repositories.CreateWorkoutPlanTx(tx, &plan); err != nil {
		tx.Rollback()
		return err
	}

	splitFocuses := helpers.GetSplitFocuses(profile.SplitType, profile.Frequency)

	restMap := map[int]bool{}
	for _, d := range restDays {
		restMap[d] = true
	}

	focusIndex := 0
	for dayNum := 1; dayNum <= 7; dayNum++ {
		if restMap[dayNum] {
			day := models.WorkoutPlanDay{
				PlanID:    plan.ID,
				DayNumber: dayNum,
				Focus:     "Rest",
			}
			if err := repositories.CreateWorkoutPlanDayTx(tx, &day); err != nil {
				tx.Rollback()
				return err
			}

			restExercise := models.WorkoutPlanExercise{
				DayID:      day.ID,
				ExerciseID: 0,
				Order:      0,
				Reps:       0,
				Sets:       0,
			}
			if err := repositories.CreateWorkoutPlanExerciseTx(tx, &restExercise); err != nil {
				tx.Rollback()
				return err
			}

			continue
		}

		if focusIndex >= len(splitFocuses) {
			break
		}
		focus := splitFocuses[focusIndex]
		focusIndex++

		day := models.WorkoutPlanDay{
			PlanID:    plan.ID,
			DayNumber: dayNum,
			Focus:     focus,
		}
		if err := repositories.CreateWorkoutPlanDayTx(tx, &day); err != nil {
			tx.Rollback()
			return err
		}

		focused := helpers.FilterExercisesByFocus(exercises, focus)
		if len(focused) == 0 {
			focused = exercises
		}

		selected := helpers.SelectTailoredExercises(focused, profile, focus, 4)
		if len(selected) == 0 {
			altProfile := *profile
			altProfile.Goal = "General Fitness"
			selected = helpers.SelectTailoredExercises(focused, &altProfile, focus, 4)
		}
		if len(selected) == 0 && strings.ToLower(profile.Intensity) == "beginner" {
			altProfile := *profile
			altProfile.Intensity = "Intermediate"
			selected = helpers.SelectTailoredExercises(focused, &altProfile, focus, 4)
		}
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
		if len(selected) == 0 {
			tx.Rollback()
			return fmt.Errorf("no suitable exercises found for focus %s", focus)
		}

		for j, ex := range selected {
			reps := helpers.DetermineReps(profile.Intensity, profile.Goal, profile.BMICategory)
			planExercise := models.WorkoutPlanExercise{
				DayID:      day.ID,
				ExerciseID: ex.ID,
				Order:      j + 1,
				Reps:       reps,
				Sets:       3,
			}
			if err := repositories.CreateWorkoutPlanExerciseTx(tx, &planExercise); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PlanService) GetAllPlans(userID uint64) ([]models.WorkoutPlan, error) {
	return repositories.GetAllWorkoutPlansByUserID(userID)
}

func (s *PlanService) GetPlanByUserID(userID uint64) (dto.FullPlanOutput, error) {
	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return dto.FullPlanOutput{}, fmt.Errorf("failed to retrieve workout plan: %w", err)
	}

	tx := config.DB.Begin()
	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return dto.FullPlanOutput{}, fmt.Errorf("failed to retrieve profile: %w", err)
	}
	tx.Commit()

	var workoutDays []dto.WorkoutDay

	for _, day := range plan.Days {
		var dayDTO dto.WorkoutDay
		dayDTO.DayNumber = day.DayNumber
		dayDTO.Focus = day.Focus

		for _, ex := range day.Exercises {
			var (
				exerciseName string
				imageURL     string
			)

			if ex.ExerciseID == 0 {
				exerciseName = "Rest Day"
				imageURL = "-"
			} else {
				exerciseDetail, err := repositories.GetExerciseByID(ex.ExerciseID)
				if err != nil {
					continue
				}
				exerciseName = exerciseDetail.Name

				blobName := "picture/image_" + fmt.Sprint(ex.ExerciseID) + ".jpg"
				imageURL, err = helpers.GenerateSASURL(blobName, time.Hour)
				if err != nil {
					imageURL = ""
				}
			}

			dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
				ExerciseID: ex.ExerciseID,
				Name:       exerciseName,
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
	tx := config.DB.Begin()
	err := repositories.DeleteWorkoutPlanByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		return err
	}
	return nil
}

func (s *PlanService) GetRecommendedReplacements(userID uint64) ([]dto.ExerciseReplacementResponse, error) {
	profile, err := repositories.GetProfileByUserID(config.DB, userID)
	if err != nil {
		return nil, err
	}
	equipment := helpers.NormalizeEquipment(profile.EquipmentJSON)

	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return nil, err
	}

	existingExerciseIDs := make(map[uint64]bool)
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			existingExerciseIDs[ex.ExerciseID] = true
		}
	}

	var results []dto.ExerciseReplacementResponse
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			exDetail, err := repositories.GetExerciseByID(ex.ExerciseID)
			if err != nil {
				continue
			}

			similar, err := repositories.FindSimilarExercises(*exDetail, profile, equipment, 10)
			if err != nil {
				continue
			}

			var recs []dto.RecommendedExerciseBrief
			for _, s := range similar {
				if s.ID == ex.ExerciseID || existingExerciseIDs[s.ID] {
					continue
				}
				recs = append(recs, dto.RecommendedExerciseBrief{
					ExerciseID:  s.ID,
					Name:        s.Name,
					Description: s.Description,
				})
				if len(recs) >= 3 {
					break
				}
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

func (s *PlanService) ReplaceExercise(userID uint64, req dto.ReplaceExerciseRequest) error {
	tx := config.DB.Begin()

	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		tx.Rollback()
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
		tx.Rollback()
		return fmt.Errorf("original exercise not found in your plan")
	}
	err = repositories.UpdateExerciseInPlanExercise(tx, targetPlanExerciseID, req.NewExerciseID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update exercise: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PlanService) EditReps(userID uint64, input dto.EditRepsRequest) error {
	tx := config.DB.Begin()

	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("user has no active workout plan")
	}

	var targetPlanExerciseID uint64
	found := false
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			if ex.ExerciseID == input.PlanExerciseID {
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
		tx.Rollback()
		return fmt.Errorf("exercise not found in your plan")
	}
	err = repositories.UpdateWorkoutPlanExerciseReps(tx, targetPlanExerciseID, input.NewReps)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update reps for exercise: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *PlanService) GetWorkoutToday(userID uint64, dayID uint64) (dto.FullDayPlanOutput, error) {
	tx := config.DB.Begin()

	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		tx.Rollback()
		return dto.FullDayPlanOutput{}, fmt.Errorf("user has no active workout plan")
	}

	var day models.WorkoutPlanDay
	err = tx.Where("plan_id = ? AND id = ?", plan.ID, dayID).First(&day).Error
	if err != nil {
		tx.Rollback()
		return dto.FullDayPlanOutput{}, fmt.Errorf("workout plan day not found")
	}

	var exercises []models.WorkoutPlanExercise
	err = tx.Where("day_id = ?", day.ID).Find(&exercises).Error
	if err != nil {
		tx.Rollback()
		return dto.FullDayPlanOutput{}, fmt.Errorf("failed to fetch exercises for the day")
	}

	var workoutDayOutput dto.WorkoutDay
	workoutDayOutput.DayNumber = day.DayNumber
	workoutDayOutput.Focus = day.Focus

	var allGoalTags []string

	for _, ex := range exercises {
		exerciseDetail, err := repositories.GetExerciseByID(ex.ExerciseID)
		if err != nil {
			continue
		}

		blobName := "picture/image_" + fmt.Sprint(ex.ExerciseID) + ".jpg"
		imageURL, err := helpers.GenerateSASURL(blobName, time.Hour)
		if err != nil {
			imageURL = ""
		}

		workoutDayOutput.Exercises = append(workoutDayOutput.Exercises, dto.ExercisePlanResponse{
			ExerciseID: ex.ExerciseID,
			Name:       exerciseDetail.Name,
			Reps:       ex.Reps,
			Sets:       ex.Sets,
			Order:      ex.Order,
			ImageURL:   imageURL,
		})

		allGoalTags = append(allGoalTags, exerciseDetail.GoalTag)
	}

	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return dto.FullDayPlanOutput{}, fmt.Errorf("failed to retrieve profile: %w", err)
	}

	output := dto.FullDayPlanOutput{
		WorkoutDay:     workoutDayOutput,
		CaloriesBurned: helpers.CalculateTodayCalories(allGoalTags, profile.TargetWeight),
	}

	if err := tx.Commit().Error; err != nil {
		return dto.FullDayPlanOutput{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return output, nil
}
