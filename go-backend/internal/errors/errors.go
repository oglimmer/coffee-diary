// Migrated from: GlobalExceptionHandler.java
package errors

import (
	"encoding/json"
	"net/http"
)

type AppError struct {
	Status  int    `json:"status"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

type ValidationError struct {
	Status  int               `json:"status"`
	Err     string            `json:"error"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

func BadRequest(message string) *AppError {
	return &AppError{Status: 400, Err: "Bad Request", Message: message}
}

func Unauthorized(message string) *AppError {
	return &AppError{Status: 401, Err: "Unauthorized", Message: message}
}

func Forbidden(message string) *AppError {
	return &AppError{Status: 403, Err: "Forbidden", Message: message}
}

func InternalError() *AppError {
	return &AppError{Status: 500, Err: "Internal Server Error", Message: "An unexpected error occurred"}
}

func WriteError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func WriteValidationError(w http.ResponseWriter, fieldErrors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(ValidationError{
		Status:  400,
		Err:     "Validation Failed",
		Message: "Invalid request body",
		Errors:  fieldErrors,
	})
}
