package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"taskmanager/internal/domain"
	"taskmanager/internal/utils"
)

// Context keys under which the authenticated principal is stored.
const (
	ctxUserID = "auth.userID"
	ctxRole   = "auth.role"
)

// Auth returns middleware that validates the Bearer token and stores the
// authenticated user's ID and role in the request context. Requests without a
// valid token are rejected with 401.
func Auth(jwt *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := extractBearer(header)
		if token == "" {
			utils.Fail(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwt.Verify(token)
		if err != nil {
			utils.Fail(c, utils.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

func extractBearer(header string) string {
	const prefix = "Bearer "
	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return strings.TrimSpace(header[len(prefix):])
	}
	return ""
}

// UserID returns the authenticated user's ID from the context.
func UserID(c *gin.Context) uuid.UUID {
	if v, ok := c.Get(ctxUserID); ok {
		if id, ok := v.(uuid.UUID); ok {
			return id
		}
	}
	return uuid.Nil
}

// CurrentRole returns the authenticated user's role from the context.
func CurrentRole(c *gin.Context) domain.Role {
	if v, ok := c.Get(ctxRole); ok {
		if r, ok := v.(domain.Role); ok {
			return r
		}
	}
	return domain.RoleUser
}
