package controllers

import (
	"strconv"
	"wellnesspath/helpers"
	"wellnesspath/services"

	"github.com/gin-gonic/gin"
)

func GetAllExercises(c *gin.Context) {
	exercises, err := (&services.ExerciseService{}).GetAllExercises()
	if err != nil {
		errorRes, status := helpers.GetErrorResponse(err)
		c.JSON(status, errorRes)
		return
	}
	helpers.SuccessResponseWithData(c, "exercises retrieved successfully", exercises)
}

func GetExerciseByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid ID", "ID must be a valid number")
		return
	}

	exercise, err := (&services.ExerciseService{}).GetExerciseByID(id)
	if err != nil {
		errorRes, status := helpers.GetErrorResponse(err)
		c.JSON(status, errorRes)
		return
	}

	helpers.SuccessResponseWithData(c, "exercise retrieved successfully", exercise)
}
