package dto

import (
	"time"

	"github.com/google/uuid"

	"taskmanager/internal/domain"
)

// --- Auth ---

// SignupRequest is the payload for account creation.
type SignupRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=128"`
}

// LoginRequest is the payload for authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is returned after a successful signup or login.
type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// --- Tasks ---

// CreateTaskRequest is the payload for creating a task.
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=200"`
	Description string     `json:"description" binding:"max=5000"`
	Status      string     `json:"status" binding:"omitempty,oneof=TODO IN_PROGRESS COMPLETED"`
	Priority    string     `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	DueDate     *time.Time `json:"dueDate" binding:"omitempty"`
}

// UpdateTaskRequest is the payload for updating a task. All fields are optional;
// pointers distinguish "absent" from "set to zero value".
type UpdateTaskRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=200"`
	Description *string    `json:"description" binding:"omitempty,max=5000"`
	Status      *string    `json:"status" binding:"omitempty,oneof=TODO IN_PROGRESS COMPLETED"`
	Priority    *string    `json:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	DueDate     *time.Time `json:"dueDate" binding:"omitempty"`
	ClearDue    bool       `json:"clearDueDate"`
}

// TaskQuery captures filtering, search, sort, and pagination parameters.
type TaskQuery struct {
	Search   string `form:"search"`
	Status   string `form:"status" binding:"omitempty,oneof=TODO IN_PROGRESS COMPLETED"`
	Priority string `form:"priority" binding:"omitempty,oneof=LOW MEDIUM HIGH"`
	Sort     string `form:"sort" binding:"omitempty,oneof=due_date_asc due_date_desc priority_asc priority_desc created_asc created_desc title_asc title_desc"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// Normalize applies defaults for unset pagination/sort values.
func (q *TaskQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Sort == "" {
		q.Sort = "created_desc"
	}
}

// Offset returns the SQL OFFSET implied by the page/limit.
func (q *TaskQuery) Offset() int {
	return (q.Page - 1) * q.Limit
}

// PageMeta is the pagination metadata returned with list responses.
type PageMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"totalPages"`
}

// NewPageMeta computes pagination metadata.
func NewPageMeta(total int64, page, limit int) PageMeta {
	totalPages := int(0)
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}
	return PageMeta{Total: total, Page: page, Limit: limit, TotalPages: totalPages}
}

// --- Admin ---

// ActivityLogQuery captures pagination for the admin activity log endpoint.
type ActivityLogQuery struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// Normalize applies defaults.
func (q *ActivityLogQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
}

// Offset returns the implied SQL OFFSET.
func (q *ActivityLogQuery) Offset() int { return (q.Page - 1) * q.Limit }

// --- SSE ---

// TaskEvent is the payload broadcast over SSE when a task changes.
type TaskEvent struct {
	Type   string       `json:"type"` // task.created | task.updated | task.deleted
	TaskID uuid.UUID    `json:"taskId"`
	Task   *domain.Task `json:"task,omitempty"`
}
