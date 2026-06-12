package handler

import (
	"github.com/gin-gonic/gin"

	"taskmanager/internal/dto"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

// AdminHandler exposes admin-only endpoints.
type AdminHandler struct {
	admin    *service.AdminService
	activity *service.ActivityService
}

// NewAdminHandler constructs an AdminHandler.
func NewAdminHandler(admin *service.AdminService, activity *service.ActivityService) *AdminHandler {
	return &AdminHandler{admin: admin, activity: activity}
}

// ListUsers handles GET /api/admin/users.
func (h *AdminHandler) ListUsers(c *gin.Context) {
	var q dto.ActivityLogQuery
	if !bindQuery(c, &q) {
		return
	}
	users, meta, err := h.admin.ListUsers(q)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Paginated(c, users, meta)
}

// ListActivityLogs handles GET /api/admin/activity-logs.
func (h *AdminHandler) ListActivityLogs(c *gin.Context) {
	var q dto.ActivityLogQuery
	if !bindQuery(c, &q) {
		return
	}
	logs, meta, err := h.activity.List(q)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Paginated(c, logs, meta)
}
