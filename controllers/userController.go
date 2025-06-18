package controllers

import (
	"wellnesspath/dto"
	"wellnesspath/helpers"
	"wellnesspath/services"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, helpers.NewUnauthorizedError("User ID not found in context"))
		return
	}

	userID, ok := userIDRaw.(uint64)
	if !ok {
		helpers.ErrorResponse(c, helpers.NewBadRequestError("Invalid User ID type"))
		return
	}

	request := dto.GetUserDTO{ID: userID}

	userResponse, err := (&services.UserService{}).GetUserByID(request)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "User fetched successfully", userResponse)
}

func UpdateUser(c *gin.Context) {
	var request *dto.UpdateUserDTO

	if err := c.ShouldBind(&request); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	if request.Name == "" && request.Profile == nil && request.Password == "" {
		helpers.ValidationErrorResponse(c, "At least one of Name, Profile, or Password must be provided", "")
		return
	}

	userResponse, err := (&services.UserService{}).UpdateUser(*request)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "User updated successfully", userResponse)
}

func DeleteUser(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		helpers.ErrorResponse(c, helpers.NewUnauthorizedError("User ID not found in context"))
		return
	}

	userID, ok := userIDRaw.(uint64)
	if !ok {
		helpers.ErrorResponse(c, helpers.NewBadRequestError("Invalid User ID type"))
		return
	}

	request := dto.DeleteUserDTO{ID: userID}

	if err := (&services.UserService{}).DeleteUser(request); err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponse(c, "User deleted successfully")
}
