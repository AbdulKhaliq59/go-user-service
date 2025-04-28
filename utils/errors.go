// utils/errors.go
package utils

import "net/http"

// CustomError is a type for handling HTTP errors with custom status codes
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewError creates a new CustomError with the given status code and message
func NewError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// ErrorResponse creates a response map for error handling
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"error": message,
	}
}

// NewBadRequestError returns a bad request error
func NewBadRequestError(message string) *CustomError {
	return NewError(http.StatusBadRequest, message)
}

// NewUnauthorizedError returns an unauthorized error
func NewUnauthorizedError(message string) *CustomError {
	return NewError(http.StatusUnauthorized, message)
}

// NewNotFoundError returns a not found error
func NewNotFoundError(message string) *CustomError {
	return NewError(http.StatusNotFound, message)
}

// NewConflictError returns a conflict error
func NewConflictError(message string) *CustomError {
	return NewError(http.StatusConflict, message)
}

// NewInternalServerError returns an internal server error
func NewInternalServerError(message string) *CustomError {
	return NewError(http.StatusInternalServerError, message)
}
