package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e *AppError) Error() string {
	return e.Message
}

// Domain Errors
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    "validation_error",
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Type:    "not_found_error",
	}
}

func NewUnauthorizedError() *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized access",
		Type:    "unauthorized_error",
	}
}

func NewForbiddenError() *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "forbidden access",
		Type:    "forbidden_error",
	}
}

func NewInternalError() *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Type:    "internal_error",
	}
}
