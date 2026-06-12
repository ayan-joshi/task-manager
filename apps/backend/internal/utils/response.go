package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Envelope is the consistent JSON shape returned by every endpoint. Exactly one
// of Data or Error is populated.
type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

// ErrorBody is the error portion of an Envelope.
type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// OK writes a 200 success response.
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

// Created writes a 201 success response.
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

// Paginated writes a 200 success response with pagination metadata.
func Paginated(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data, Meta: meta})
}

// Fail writes an error response, translating AppError into its status/code and
// falling back to 500 for unexpected errors.
func Fail(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Status, Envelope{
			Success: false,
			Error:   &ErrorBody{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Envelope{
		Success: false,
		Error:   &ErrorBody{Code: ErrInternal.Code, Message: ErrInternal.Message},
	})
}

// FailValidation writes a 422 response with field-level validation details.
func FailValidation(c *gin.Context, details interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "request validation failed",
			Details: details,
		},
	})
}
