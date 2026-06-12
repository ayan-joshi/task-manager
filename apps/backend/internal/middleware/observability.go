package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"taskmanager/internal/utils"
)

// RequestID attaches a unique request identifier to each request and echoes it
// in the X-Request-ID response header, aiding log correlation.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}
		c.Set("requestID", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

// Logger emits a structured single-line log per request.
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		c.Next()
		log.Printf("%s %s %d %s req=%v",
			c.Request.Method, path, c.Writer.Status(),
			time.Since(start), c.GetString("requestID"))
	}
}

// Recovery converts panics into a clean 500 response instead of crashing the
// process, and logs the failure with its request ID.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("panic recovered req=%v: %v", c.GetString("requestID"), r)
				utils.Fail(c, utils.ErrInternal)
				c.Abort()
			}
		}()
		c.Next()
	}
}
