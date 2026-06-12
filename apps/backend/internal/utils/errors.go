package utils

import "net/http"

// AppError is a domain-level error carrying an HTTP status code and a stable,
// machine-readable code alongside a human-readable message. Services return
// these so handlers can translate them into consistent HTTP responses without
// leaking internal details.
type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string { return e.Message }

// NewAppError builds an AppError.
func NewAppError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

// Common, reusable application errors.
var (
	ErrEmailTaken      = NewAppError(http.StatusConflict, "EMAIL_TAKEN", "email is already registered")
	ErrInvalidCreds    = NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid email or password")
	ErrUnauthorized    = NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
	ErrForbidden       = NewAppError(http.StatusForbidden, "FORBIDDEN", "you do not have permission to perform this action")
	ErrTaskNotFound    = NewAppError(http.StatusNotFound, "TASK_NOT_FOUND", "task not found")
	ErrAttachNotFound  = NewAppError(http.StatusNotFound, "ATTACHMENT_NOT_FOUND", "attachment not found")
	ErrFileTooLarge    = NewAppError(http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "uploaded file exceeds the maximum allowed size")
	ErrUnsupportedType = NewAppError(http.StatusUnsupportedMediaType, "UNSUPPORTED_TYPE", "uploaded file type is not allowed")
	ErrInternal        = NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "an unexpected error occurred")
)
