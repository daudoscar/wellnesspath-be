package controllers

import (
	"errors"
	"strconv"
	"wellnesspath/helpers"
	"wellnesspath/services"

	"github.com/gin-gonic/gin"
)

func GetAds(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	adsIDStr := c.Query("adsID")
	if adsIDStr == "" {
		helpers.ValidationErrorResponse(c, "adsID is required", "")
		return
	}

	adsID, err := strconv.ParseUint(adsIDStr, 10, 64)
	if err != nil {
		helpers.ValidationErrorResponse(c, "Invalid adsID format", err.Error())
		return
	}

	AdResponse, err := (&services.AdsService{}).GetAds(userID.(uint64), adsID)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Workout for today fetched successfully", AdResponse)
}
