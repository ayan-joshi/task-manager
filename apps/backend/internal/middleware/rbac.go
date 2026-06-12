package middleware

import (
	"github.com/gin-gonic/gin"

	"taskmanager/internal/domain"
	"taskmanager/internal/utils"
)

// RequireRole returns middleware that allows the request only if the
// authenticated user holds one of the permitted roles. It must run after Auth.
func RequireRole(roles ...domain.Role) gin.HandlerFunc {
	allowed := make(map[domain.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		if _, ok := allowed[CurrentRole(c)]; !ok {
			utils.Fail(c, utils.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
