package services

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if err != nil || len(exercises) == 0 {
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

	usedExerciseIDs := map[uint64]bool{}
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
			focused = exercises // fallback ke semua
		}

		reps := helpers.DetermineReps(profile.Intensity, profile.Goal, profile.BMICategory)
		exerciseCount := helpers.CalculateMaxExercises(profile.DurationPerSession, reps)
		validParts := helpers.GetBodyPartsForFocus(focus)

		selected := []models.Exercise{}
		candidate := helpers.FilterWithBodyPartCoverage(focused, validParts, profile.Goal, profile.Intensity, exerciseCount)

		for _, ex := range candidate {
			if !usedExerciseIDs[ex.ID] {
				selected = append(selected, ex)
				usedExerciseIDs[ex.ID] = true
			}
			if len(selected) == exerciseCount {
				break
			}
		}

		// Fallback jika belum cukup
		if len(selected) < exerciseCount {
			for _, ex := range focused {
				if !usedExerciseIDs[ex.ID] {
					selected = append(selected, ex)
					usedExerciseIDs[ex.ID] = true
				}
				if len(selected) == exerciseCount {
					break
				}
			}
		}

		if len(selected) == 0 {
			tx.Rollback()
			return fmt.Errorf("no suitable exercises found for focus %s", focus)
		}

		for j, ex := range selected {
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

func (s *PlanService) InitializeWorkoutPlan(userID uint64) (*dto.CreateDaysRequest, error) {
	tx := config.DB.Begin()

	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user profile not found")
	}

	var restDays []int
	if err := json.Unmarshal([]byte(profile.RestDaysJSON), &restDays); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to parse rest days: %w", err)
	}

	if err := helpers.ValidateSplitAndRestDays(profile.SplitType, profile.Frequency, restDays); err != nil {
		tx.Rollback()
		return nil, err
	}

	existingPlans, err := repositories.GetAllWorkoutPlansByUserID(userID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to check existing workout plans: %w", err)
	}
	if len(existingPlans) > 0 {
		if err := repositories.DeleteWorkoutPlanByUserID(tx, userID); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete existing workout plans: %w", err)
		}
	}

	equipment := helpers.DecodeEquipment(profile.EquipmentJSON)
	exercises, err := repositories.GetExercisesByGoalAndEquipment(profile.Goal, equipment)
	if err != nil || len(exercises) == 0 {
		tx.Rollback()
		return nil, errors.New("no exercises match your profile")
	}

	plan := &models.WorkoutPlan{
		UserID:    userID,
		SplitType: profile.SplitType,
		Goal:      profile.Goal,
	}
	if err := repositories.CreateWorkoutPlanTx(tx, plan); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// RETURN SESUAI DTO CreateDaysRequest YANG DIBUTUHKAN FUNCTION KE-2
	return &dto.CreateDaysRequest{
		PlanID:   plan.ID,
		RestDays: restDays,
		Profile:  *profile, // seluruh profil disertakan
	}, nil
}

func (s *PlanService) CreateWorkoutPlanDays(input dto.CreateDaysRequest) (*dto.InsertExercisesRequest, error) {
	tx := config.DB.Begin()

	restMap := map[int]bool{}
	for _, d := range input.RestDays {
		restMap[d] = true
	}

	splitFocuses := helpers.GetSplitFocuses(input.Profile.SplitType, input.Profile.Frequency)
	focusIndex := 0

	var allDays []models.WorkoutPlanDay

	for dayNum := 1; dayNum <= 7; dayNum++ {
		day := models.WorkoutPlanDay{
			PlanID:    input.PlanID,
			DayNumber: dayNum,
			Focus:     "Rest",
		}

		// Assign focus jika bukan hari istirahat
		if !restMap[dayNum] {
			if focusIndex < len(splitFocuses) {
				day.Focus = splitFocuses[focusIndex]
				focusIndex++
			}
		}

		// Simpan ke DB
		if err := repositories.CreateWorkoutPlanDayTx(tx, &day); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Kalau hari rest, tambahkan exercise dummy
		if day.Focus == "Rest" {
			restExercise := models.WorkoutPlanExercise{
				DayID:      day.ID,
				ExerciseID: 0,
				Order:      0,
				Reps:       0,
				Sets:       0,
			}
			if err := repositories.CreateWorkoutPlanExerciseTx(tx, &restExercise); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		allDays = append(allDays, day)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Return format input untuk InsertExercisesToDays
	return &dto.InsertExercisesRequest{
		Profile: input.Profile,
		Days:    allDays,
	}, nil
}

func (s *PlanService) InsertExercisesToDays(input dto.InsertExercisesRequest) error {
	tx := config.DB.Begin()

	// Ambil ulang daftar exercise dari DB berdasarkan profile
	equipment := helpers.DecodeEquipment(input.Profile.EquipmentJSON)
	exercises, err := repositories.GetExercisesByGoalAndEquipment(input.Profile.Goal, equipment)
	if err != nil || len(exercises) == 0 {
		tx.Rollback()
		return errors.New("no exercises match your profile")
	}

	splitFocuses := helpers.GetSplitFocuses(input.Profile.SplitType, input.Profile.Frequency)
	usedExerciseIDs := map[uint64]bool{}
	focusIndex := 0

	for _, day := range input.Days {
		if day.Focus == "Rest" {
			continue
		}

		if focusIndex >= len(splitFocuses) {
			break
		}
		focus := splitFocuses[focusIndex]
		focusIndex++

		reps := helpers.DetermineReps(input.Profile.Intensity, input.Profile.Goal, input.Profile.BMICategory)
		exerciseCount := helpers.CalculateMaxExercises(input.Profile.DurationPerSession, reps)
		validParts := helpers.GetBodyPartsForFocus(focus)

		focused := helpers.FilterExercisesByFocus(exercises, focus)
		if len(focused) == 0 {
			focused = exercises
		}

		candidate := helpers.FilterWithBodyPartCoverage(focused, validParts, input.Profile.Goal, input.Profile.Intensity, exerciseCount)

		selected := []models.Exercise{}
		for _, ex := range candidate {
			if !usedExerciseIDs[ex.ID] {
				selected = append(selected, ex)
				usedExerciseIDs[ex.ID] = true
			}
			if len(selected) == exerciseCount {
				break
			}
		}

		if len(selected) < exerciseCount {
			for _, ex := range focused {
				if !usedExerciseIDs[ex.ID] {
					selected = append(selected, ex)
					usedExerciseIDs[ex.ID] = true
				}
				if len(selected) == exerciseCount {
					break
				}
			}
		}

		if len(selected) == 0 {
			tx.Rollback()
			return fmt.Errorf("no suitable exercises found for day %d", day.DayNumber)
		}

		// ✅ Batch Insert for the current day
		batch := make([]models.WorkoutPlanExercise, 0, len(selected))
		for i, ex := range selected {
			batch = append(batch, models.WorkoutPlanExercise{
				DayID:      day.ID,
				ExerciseID: ex.ID,
				Order:      i + 1,
				Reps:       reps,
				Sets:       3,
			})
		}

		if err := repositories.CreateWorkoutPlanExercisesBatchTx(tx, batch); err != nil {
			tx.Rollback()
			return fmt.Errorf("batch insert failed for day %d: %w", day.DayNumber, err)
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

	profile, err := repositories.GetProfileByUserID(config.DB, userID)
	if err != nil {
		return dto.FullPlanOutput{}, fmt.Errorf("failed to retrieve profile: %w", err)
	}

	// Kumpulkan semua ExerciseID
	uniqueIDs := make(map[uint64]struct{})
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			if ex.ExerciseID != 0 {
				uniqueIDs[ex.ExerciseID] = struct{}{}
			}
		}
	}

	var ids []uint64
	for id := range uniqueIDs {
		ids = append(ids, id)
	}

	exDetails, err := repositories.GetExercisesByIDs(ids)
	if err != nil {
		return dto.FullPlanOutput{}, err
	}

	exMap := make(map[uint64]*models.Exercise)
	for _, e := range exDetails {
		exMap[e.ID] = e
	}

	var workoutDays []dto.WorkoutDay
	for _, day := range plan.Days {
		var dayDTO dto.WorkoutDay
		dayDTO.DayNumber = day.DayNumber
		dayDTO.Focus = day.Focus

		for _, ex := range day.Exercises {
			if ex.ExerciseID == 0 {
				dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
					ExerciseID: 0,
					Name:       "Rest Day",
					Reps:       ex.Reps,
					Sets:       ex.Sets,
					Order:      ex.Order,
					BodyPart:   "-",
					Equipment:  "-",
				})
				continue
			}

			detail := exMap[ex.ExerciseID]
			dayDTO.Exercises = append(dayDTO.Exercises, dto.ExercisePlanResponse{
				ExerciseID: ex.ExerciseID,
				Name:       detail.Name,
				Reps:       ex.Reps,
				Sets:       ex.Sets,
				Order:      ex.Order,
				BodyPart:   detail.BodyPart,
				Equipment:  detail.Equipment,
			})
		}
		workoutDays = append(workoutDays, dayDTO)
	}

	return dto.FullPlanOutput{
		WorkoutPlan:    workoutDays,
		TrainingAdvice: helpers.GenerateTrainingAdvice(profile),
		BMIInfo:        helpers.BuildBMIInfo(profile.BMI, profile.BMICategory),
		CaloriesBurned: helpers.CalculateCalories(profile),
		NutritionPlan:  helpers.GenerateNutrition(profile),
	}, nil
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
	// 1. Ambil profil user
	profile, err := repositories.GetProfileByUserID(config.DB, userID)
	if err != nil {
		return nil, err
	}
	equipment := helpers.NormalizeEquipment(profile.EquipmentJSON)

	// 2. Ambil active workout plan
	plan, err := repositories.GetActiveWorkoutPlanByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 3. Kumpulkan semua ExerciseID di plan dan bodypart unik
	var allExerciseIDs []uint64
	existingExerciseIDs := make(map[uint64]bool)
	bodyParts := map[string]struct{}{}
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			allExerciseIDs = append(allExerciseIDs, ex.ExerciseID)
			existingExerciseIDs[ex.ExerciseID] = true
		}
	}

	// Ambil detail exercise di plan (batch)
	exerciseMap, err := repositories.GetExercisesByIDs(allExerciseIDs)
	if err != nil {
		return nil, err
	}

	// Cari bodypart unik
	for _, ex := range exerciseMap {
		bodyParts[ex.BodyPart] = struct{}{}
	}
	var uniqueBodyParts []string
	for bp := range bodyParts {
		uniqueBodyParts = append(uniqueBodyParts, bp)
	}

	// 4. Ambil kandidat replacement dengan batch query di repositories
	candidateExercises, err := repositories.FindExercisesByBodyPartsAndEquipment(config.DB, uniqueBodyParts, equipment, allExerciseIDs)
	if err != nil {
		return nil, err
	}
	// Group candidate by bodypart
	candidatesByBodyPart := map[string][]models.Exercise{}
	for _, c := range candidateExercises {
		candidatesByBodyPart[c.BodyPart] = append(candidatesByBodyPart[c.BodyPart], c)
	}

	// 5. Compose response (filter & limit di sini)
	var results []dto.ExerciseReplacementResponse
	for _, day := range plan.Days {
		for _, ex := range day.Exercises {
			exDetail := exerciseMap[ex.ExerciseID]
			candidates := candidatesByBodyPart[exDetail.BodyPart]
			var recs []dto.RecommendedExerciseBrief
			for _, s := range candidates {
				if s.ID == ex.ExerciseID {
					continue
				}
				recs = append(recs, dto.RecommendedExerciseBrief{
					ExerciseID:  s.ID,
					Name:        s.Name,
					Description: s.Description,
				})
				if len(recs) >= 3 { // limit kandidat replacement
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
		return dto.FullDayPlanOutput{}, fmt.Errorf("user has no active workout plan")
	}

	var day models.WorkoutPlanDay
	if err := tx.Where("plan_id = ? AND day_number = ?", plan.ID, dayID).First(&day).Error; err != nil {
		return dto.FullDayPlanOutput{}, fmt.Errorf("workout plan day not found")
	}

	var exercises []models.WorkoutPlanExercise
	if err := tx.Where("day_id = ?", day.ID).Find(&exercises).Error; err != nil {
		return dto.FullDayPlanOutput{}, fmt.Errorf("failed to fetch exercises for the day")
	}

	exerciseIDs := make([]uint64, 0, len(exercises))
	for _, ex := range exercises {
		exerciseIDs = append(exerciseIDs, ex.ExerciseID)
	}

	exMap, err := repositories.GetExercisesByIDs(exerciseIDs)
	if err != nil {
		return dto.FullDayPlanOutput{}, fmt.Errorf("failed to retrieve exercise details: %w", err)
	}

	var workoutDayOutput dto.WorkoutDayToday
	workoutDayOutput.DayNumber = day.DayNumber
	workoutDayOutput.Focus = day.Focus

	var allGoalTags []string
	for _, ex := range exercises {
		detail, ok := exMap[ex.ExerciseID]
		if !ok {
			continue
		}

		blobName := "picture/image_" + fmt.Sprint(ex.ExerciseID) + ".jpg"
		imageURL, err := helpers.GenerateSASURL(blobName, time.Hour)
		if err != nil {
			imageURL = ""
		}

		workoutDayOutput.Exercises = append(workoutDayOutput.Exercises, dto.ExerciseTodayResponse{
			ExerciseID: ex.ExerciseID,
			Name:       detail.Name,
			Reps:       ex.Reps,
			Sets:       ex.Sets,
			Order:      ex.Order,
			ImageURL:   imageURL,
		})

		allGoalTags = append(allGoalTags, detail.GoalTag)
	}

	profile, err := repositories.GetProfileByUserID(tx, userID)
	if err != nil {
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
