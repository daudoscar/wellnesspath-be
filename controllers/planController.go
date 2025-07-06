package controllers

import (
	"errors"
	"strconv"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/services"

	"github.com/gin-gonic/gin"
)

func GenerateWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	err := (&services.PlanService{}).GenerateWorkoutPlan(userID.(uint64))
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Workout plan generated successfully")
}

// SEPERATE FUNCTION ***
func InitializeWorkoutPlan(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	result, err := (&services.PlanService{}).InitializeWorkoutPlan(userID.(uint64))
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout plan initialized successfully", result)
}

func CreateWorkoutPlanDays(c *gin.Context) {
	var input dto.CreateDaysRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	days, err := (&services.PlanService{}).CreateWorkoutPlanDays(input)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout days created successfully", days)
}

func InsertExercisesToDays(c *gin.Context) {
	var input dto.InsertExercisesRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	if err := (&services.PlanService{}).InsertExercisesToDays(input); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Exercises inserted successfully")
}

// SEPERATE FUNCTION ***

func GetAllPlans(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	plans, err := (&services.PlanService{}).GetAllPlans(userID.(uint64))
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout plans retrieved successfully", plans)
}

func GetPlanByUserID(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	userID, ok := userIDRaw.(uint64)
	if !ok {
		helpers.ValidationErrorResponse(c, "Invalid ID", "User ID must be a valid number")
		return
	}

	plan, err := (&services.PlanService{}).GetPlanByUserID(userID)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout plan retrieved successfully", plan)
}

func DeletePlan(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, helpers.NewUnauthorizedError("User not authenticated"))
		return
	}

	userID, ok := userIDRaw.(uint64)
	if !ok {
		helpers.ErrorResponse(c, helpers.NewBadRequestError("Invalid User ID type"))
		return
	}

	if err := (&services.PlanService{}).DeletePlan(userID); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Workout plan deleted successfully")
}

func GetRecommendedReplacements(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	userID := userIDRaw.(uint64)

	replacements, err := (&services.PlanService{}).GetRecommendedReplacements(userID)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Recommended replacements retrieved", replacements)
}

func ReplaceExercise(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	var req dto.ReplaceExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	err := (&services.PlanService{}).ReplaceExercise(userID.(uint64), req)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Exercise replaced successfully")
}

func UpdateExerciseReps(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	var req dto.EditRepsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request body", err.Error())
		return
	}

	err := (&services.PlanService{}).EditReps(userID.(uint64), req)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Reps updated successfully")
}

func GetWorkoutToday(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	dayIDStr := c.Query("dayID")
	if dayIDStr == "" {
		helpers.ValidationErrorResponse(c, "dayID is required", "")
		return
	}

	dayID, err := strconv.ParseUint(dayIDStr, 10, 64)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid dayID format", err.Error())
		return
	}

	plan, err := (&services.PlanService{}).GetWorkoutToday(userID.(uint64), dayID)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout for today fetched successfully", plan)
}
