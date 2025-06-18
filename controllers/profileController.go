package controllers

import (
	"errors"
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/services"

	"github.com/gin-gonic/gin"
)

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	profile, err := (&services.ProfileService{}).GetProfile(userID.(uint64))
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "profile retrieved successfully", profile)
}

func UpdateProfile(c *gin.Context) {
	var request dto.UpdateProfileDTO

	if err := c.ShouldBindJSON(&request); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	if err := (&services.ProfileService{}).UpdateProfile(userID.(uint64), request); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "profile updated successfully")
}

func DeleteProfile(c *gin.Context) {
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

	if err := (&services.ProfileService{}).DeleteProfile(userID); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "Profile deleted successfully")
}
