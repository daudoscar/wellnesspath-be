package helpers

import (
	"fmt"
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

	if strings.HasPrefix(err.Error(), "unauthorized:") {
		errorResponse.Message = strings.TrimPrefix(err.Error(), "unauthorized:")
		return errorResponse, http.StatusUnauthorized
	}

	if strings.HasPrefix(err.Error(), "bad_request:") {
		errorResponse.Message = strings.TrimPrefix(err.Error(), "bad_request:")
		return errorResponse, http.StatusBadRequest
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

func NewUnauthorizedError(message string) error {
	return fmt.Errorf("unauthorized: %s", message)
}

func NewBadRequestError(message string) error {
	return fmt.Errorf("bad_request: %s", message)
}
