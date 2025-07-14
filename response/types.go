// ==================== response/types.go ====================
package response

import (
	"fmt"
)

// ErrorCode represents application error codes
type ErrorCode string

const (
	// Client errors (4xx)
	ErrCodeInvalidInput    ErrorCode = "INVALID_INPUT"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeConflict        ErrorCode = "CONFLICT"
	ErrCodeRequestTooLarge ErrorCode = "REQUEST_TOO_LARGE"
	ErrCodeTooManyRequest  ErrorCode = "TOO_MANY_REQUESTS"

	// Server errors (5xx)
	ErrCodeInternalServer  ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeDatabaseError   ErrorCode = "DATABASE_ERROR"
	ErrCodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// AppError represents application error with context
type AppError struct {
	Code       ErrorCode      `json:"code"`
	Message    string         `json:"message"`
	HTTPStatus int            `json:"-"`
	Err        error          `json:"-"`
	Context    map[string]any `json:"context,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) WithContext(key string, value any) *AppError {
	if e.Context == nil {
		e.Context = make(map[string]any)
	}
	e.Context[key] = value
	return e
}

// ErrorResponse represents error response structure
type ErrorResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Code    ErrorCode      `json:"code"`
	Errors  map[string]any `json:"errors,omitempty"`
}

// SuccessResponse represents success response structure
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

// Meta represents metadata for responses
type Meta struct {
	Pagination  *Pagination     `json:"pagination,omitempty"`
	Permissions map[string]bool `json:"permissions,omitempty"`
	Flags       map[string]bool `json:"flags,omitempty"`
}

// Pagination represents pagination information
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
	Offset     int `json:"offset"`
}
