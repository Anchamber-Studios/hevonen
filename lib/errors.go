package lib

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrLoginFailed   = errors.New("login failed")
	ErrUnauthorized  = errors.New("unauthorized")
)

type ApiError struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"errorCode"`
	Message    string `json:"message"`
}

func (e *ApiError) Error() string {
	return "api error"
}

func NewNotFoundError() *ApiError {
	return &ApiError{
		StatusCode: http.StatusNotFound,
		ErrorCode:  "not_found",
		Message:    "entity not found",
	}
}

func NewUnauthorizedError() *ApiError {
	return &ApiError{
		StatusCode: http.StatusUnauthorized,
		ErrorCode:  "unauthorized",
		Message:    "missing authorizations for this action",
	}
}

func NewUnauthendicatedError() *ApiError {
	return &ApiError{
		StatusCode: http.StatusForbidden,
		ErrorCode:  "unauthenticated",
		Message:    "user is not authenticated",
	}
}

func (e *ApiError) WithMessage(message string) *ApiError {
	e.Message = message
	return e
}

func (e *ApiError) WithStatusCode(statusCode int) *ApiError {
	e.StatusCode = statusCode
	return e
}

func (e *ApiError) WithErrorCode(errorCode string) *ApiError {
	e.ErrorCode = errorCode
	return e
}

type ValidationError struct {
	ErrorCode      string
	Message        string
	TranslationKey string
	Field          string
	Children       FieldErrors
}

type FieldErrors map[string]*ValidationError

func NewValidationError() *ValidationError {
	return &ValidationError{
		ErrorCode: "validation_error",
		Message:   "validation error",
		Field:     "",
		Children:  map[string]*ValidationError{},
	}
}
func NewValidationErrorForFields(errors FieldErrors) *ValidationError {
	return &ValidationError{
		ErrorCode: "validation_error",
		Message:   "validation error",
		Children:  errors,
	}
}

func (e *ValidationError) Error() string {
	return "api error"
}

func (e *ValidationError) WithMessage(message string) *ValidationError {
	e.Message = message
	return e
}

func (e *ValidationError) WithErrorCode(errorCode string) *ValidationError {
	e.ErrorCode = errorCode
	return e
}

func (e *ValidationError) WithField(field string) *ValidationError {
	e.Field = field
	return e
}
