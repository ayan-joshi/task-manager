package handler

import (
	"github.com/gin-gonic/gin"

	"taskmanager/internal/dto"
	"taskmanager/internal/middleware"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

// TaskHandler exposes task CRUD and listing endpoints.
type TaskHandler struct {
	tasks *service.TaskService
}

// NewTaskHandler constructs a TaskHandler.
func NewTaskHandler(tasks *service.TaskService) *TaskHandler {
	return &TaskHandler{tasks: tasks}
}

// Create handles POST /api/tasks.
func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Create(middleware.UserID(c), req)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Created(c, task)
}

// List handles GET /api/tasks with filtering, search, sort, and pagination.
func (h *TaskHandler) List(c *gin.Context) {
	var q dto.TaskQuery
	if !bindQuery(c, &q) {
		return
	}
	tasks, meta, err := h.tasks.List(middleware.UserID(c), middleware.CurrentRole(c), q)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Paginated(c, tasks, meta)
}

// Get handles GET /api/tasks/:id.
func (h *TaskHandler) Get(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	task, err := h.tasks.GetByID(middleware.UserID(c), middleware.CurrentRole(c), id)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, task)
}

// Update handles PUT /api/tasks/:id.
func (h *TaskHandler) Update(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	var req dto.UpdateTaskRequest
	if !bindJSON(c, &req) {
		return
	}
	task, err := h.tasks.Update(middleware.UserID(c), middleware.CurrentRole(c), id, req)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, task)
}

// Activity handles GET /api/tasks/:id/activity, returning the task's audit trail.
func (h *TaskHandler) Activity(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	var q dto.ActivityLogQuery
	if !bindQuery(c, &q) {
		return
	}
	logs, meta, err := h.tasks.Activity(middleware.UserID(c), middleware.CurrentRole(c), id, q)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Paginated(c, logs, meta)
}

// Delete handles DELETE /api/tasks/:id.
func (h *TaskHandler) Delete(c *gin.Context) {
	id, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	if err := h.tasks.Delete(middleware.UserID(c), middleware.CurrentRole(c), id); err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, gin.H{"deleted": true, "id": id})
}
