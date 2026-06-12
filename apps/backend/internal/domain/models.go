package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role enumerates the authorization roles a user can hold.
type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

// Status enumerates the lifecycle states of a task.
type Status string

const (
	StatusTodo       Status = "TODO"
	StatusInProgress Status = "IN_PROGRESS"
	StatusCompleted  Status = "COMPLETED"
)

// Priority enumerates task priority levels.
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

// Action enumerates the kinds of activity recorded in the activity log.
type Action string

const (
	ActionTaskCreated       Action = "TASK_CREATED"
	ActionTaskUpdated       Action = "TASK_UPDATED"
	ActionTaskDeleted       Action = "TASK_DELETED"
	ActionTaskStatusChanged Action = "TASK_STATUS_CHANGED"
	ActionAttachmentAdded   Action = "ATTACHMENT_ADDED"
	ActionAttachmentDeleted Action = "ATTACHMENT_DELETED"
)

// User is a registered account. The password hash is never serialized to JSON.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(120);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"column:password_hash;not null" json:"-"`
	Role      Role      `gorm:"type:varchar(20);not null;default:USER" json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	Tasks []Task `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
}

// BeforeCreate assigns a UUID primary key when one is not already set. Generating
// IDs in application code keeps the model database-agnostic (works on both
// PostgreSQL in production and SQLite in tests).
func (u *User) BeforeCreate(_ *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.Role == "" {
		u.Role = RoleUser
	}
	return nil
}

// Task is a unit of work owned by a single user.
type Task struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;index;not null" json:"userId"`
	Title       string     `gorm:"type:varchar(200);not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	Status      Status     `gorm:"type:varchar(20);index;not null;default:TODO" json:"status"`
	Priority    Priority   `gorm:"type:varchar(20);index;not null;default:MEDIUM" json:"priority"`
	DueDate     *time.Time `json:"dueDate"`
	CreatedAt   time.Time  `gorm:"index" json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`

	Attachments []Attachment `gorm:"constraint:OnDelete:CASCADE;" json:"attachments,omitempty"`
}

// BeforeCreate assigns a UUID and applies enum defaults.
func (t *Task) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	if t.Status == "" {
		t.Status = StatusTodo
	}
	if t.Priority == "" {
		t.Priority = PriorityMedium
	}
	return nil
}

// Attachment is a file associated with a task.
type Attachment struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TaskID       uuid.UUID `gorm:"type:uuid;index;not null" json:"taskId"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null" json:"userId"`
	FileName     string    `gorm:"type:varchar(255);not null" json:"fileName"`
	StoredName   string    `gorm:"type:varchar(255);not null" json:"-"`
	ContentType  string    `gorm:"type:varchar(120);not null" json:"contentType"`
	SizeBytes    int64     `json:"sizeBytes"`
	CreatedAt    time.Time `json:"createdAt"`
}

// BeforeCreate assigns a UUID.
func (a *Attachment) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// ActivityLog records an auditable action taken against a task.
type ActivityLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;index;not null" json:"userId"`
	TaskID    *uuid.UUID `gorm:"type:uuid;index" json:"taskId"`
	Action    Action     `gorm:"type:varchar(40);not null" json:"action"`
	Metadata  string     `gorm:"type:text" json:"metadata"`
	CreatedAt time.Time  `gorm:"index" json:"createdAt"`

	User *User `gorm:"constraint:OnDelete:CASCADE;" json:"user,omitempty"`
}

// BeforeCreate assigns a UUID.
func (l *ActivityLog) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
