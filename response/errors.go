// Error constructors
package response

import "net/http"

// IsAppError checks if error is AppError
func IsAppError(err error) (*AppError, bool) {
	if appErr, ok := err.(*AppError); ok {
		return appErr, true
	}
	return nil, false
}

// IsServerError checks if error is server error
func IsServerError(err error) bool {
	if appErr, ok := IsAppError(err); ok {
		return appErr.HTTPStatus >= 500
	}
	return false
}

// Error constructors
func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:       ErrCodeInvalidInput,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewForbidden(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

func NewNotFound(message string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

func NewConflict(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func NewRequestTooLarge(message string) *AppError {
	return &AppError{
		Code:       ErrCodeRequestTooLarge,
		Message:    message,
		HTTPStatus: http.StatusRequestEntityTooLarge,
	}
}

func NewTooManyRequests(message string) *AppError {
	return &AppError{
		Code:       ErrCodeTooManyRequest,
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
	}
}

func NewInternalServerError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternalServer,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewDatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeDatabaseError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewExternalServiceError(message string, err error) *AppError {
	return &AppError{
		Code:       ErrCodeExternalService,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}
