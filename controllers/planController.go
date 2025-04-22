package controllers

import (
	"errors"
	"strconv"
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

	planResponse, err := (&services.PlanService{}).GenerateWorkoutPlan(userID.(uint64))
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout plan generated successfully", planResponse)
}

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

func GetPlanByID(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid ID", "Plan ID must be a valid number")
		return
	}

	plan, err := (&services.PlanService{}).GetPlanByID(planID)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout plan retrieved successfully", plan)
}
