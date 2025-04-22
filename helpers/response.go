package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseWithData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ResponseWithoutData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func SuccessResponseWithData(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, ResponseWithData{
		Status:  "Success",
		Message: message,
		Data:    data,
	})
}

func SuccessResponse(c *gin.Context, message string) {
	c.JSON(http.StatusOK, ResponseWithoutData{
		Status:  "Success",
		Message: message,
	})
}

func ValidationErrorResponse(c *gin.Context, message string, details any) {
	c.JSON(http.StatusBadRequest, ResponseWithData{
		Status:  "Error",
		Message: message,
		Data:    details,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	errorResponse, statusCode := GetErrorResponse(err)
	c.JSON(statusCode, errorResponse)
}
