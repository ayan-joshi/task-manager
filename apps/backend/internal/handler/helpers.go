package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"taskmanager/internal/utils"
)

// fieldError is a single field-level validation failure.
type fieldError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

// bindJSON binds and validates a JSON body, writing a 422 response with
// field-level details on failure. It returns false when binding failed so the
// caller can return early.
func bindJSON(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		writeBindError(c, err)
		return false
	}
	return true
}

// bindQuery binds and validates query parameters.
func bindQuery(c *gin.Context, target interface{}) bool {
	if err := c.ShouldBindQuery(target); err != nil {
		writeBindError(c, err)
		return false
	}
	return true
}

func writeBindError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		details := make([]fieldError, 0, len(ve))
		for _, fe := range ve {
			details = append(details, fieldError{Field: fe.Field(), Rule: fe.Tag()})
		}
		utils.FailValidation(c, details)
		return
	}
	utils.FailValidation(c, err.Error())
}

// parseUUIDParam parses a path parameter as a UUID, writing a 404 on failure
// (an unparseable ID can never identify an existing resource).
func parseUUIDParam(c *gin.Context, name string) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param(name))
	if err != nil {
		utils.Fail(c, utils.ErrTaskNotFound)
		return uuid.Nil, false
	}
	return id, true
}
