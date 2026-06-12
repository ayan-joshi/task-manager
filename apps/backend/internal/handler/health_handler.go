package handler

import (
	"github.com/gin-gonic/gin"

	"taskmanager/internal/utils"
)

// Health handles GET /api/health, a lightweight liveness probe used by load
// balancers and container orchestrators.
func Health(c *gin.Context) {
	utils.OK(c, gin.H{"status": "ok"})
}
