package controllers

import (
	"syncspend/dto"
	"syncspend/helpers"
	"syncspend/services"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var request *dto.CreateUserDTO

	if err := c.ShouldBind(&request); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	credentialResponse, err := (&services.AuthService{}).RegisterUser(*request)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "User registered successfully", credentialResponse)
}

func Login(c *gin.Context) {
	var loginRequest *dto.LoginCredentialsDTO

	if err := c.ShouldBind(&loginRequest); err != nil {
		helpers.ValidationErrorResponse(c, "Invalid request", err.Error())
		return
	}

	CredentialResponseDTO, err := (&services.AuthService{}).AuthenticateUser(*loginRequest)
	if err != nil {
		helpers.ErrorResponse(c, err)
		return
	}

	helpers.SuccessResponseWithData(c, "Login successful", CredentialResponseDTO)
}
