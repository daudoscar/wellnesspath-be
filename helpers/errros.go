package helpers

import (
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type CustomError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func GetErrorResponse(err error) (CustomError, int) {
	errorResponse := CustomError{
		Status:  "Error",
		Message: "Internal Server Error",
	}

	if strings.Contains(err.Error(), "1062") {
		errorResponse.Message = "Duplicate entry error"
		return errorResponse, http.StatusConflict
	}

	if err == gorm.ErrRecordNotFound {
		errorResponse.Message = "Record not found"
		return errorResponse, http.StatusNotFound
	}

	errorResponse.Message = err.Error()

	return errorResponse, http.StatusInternalServerError
}
