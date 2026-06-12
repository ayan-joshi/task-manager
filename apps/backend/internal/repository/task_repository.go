package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"taskmanager/internal/domain"
	"taskmanager/internal/dto"
)

// TaskFilter narrows a task query. A nil UserID means "all users" (admin scope).
type TaskFilter struct {
	UserID *uuid.UUID
	Query  dto.TaskQuery
}

// TaskRepository abstracts persistence for tasks.
type TaskRepository interface {
	Create(task *domain.Task) error
	Update(task *domain.Task) error
	Delete(task *domain.Task) error
	FindByID(id uuid.UUID) (*domain.Task, error)
	List(filter TaskFilter) ([]domain.Task, int64, error)
}

type taskRepository struct{ db *gorm.DB }

// NewTaskRepository constructs a TaskRepository backed by GORM.
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) Update(task *domain.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(task *domain.Task) error {
	return r.db.Delete(task).Error
}

func (r *taskRepository) FindByID(id uuid.UUID) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Preload("Attachments").First(&task, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// List applies filtering, search, sorting, and pagination entirely in the
// database. The count and the page query share the same WHERE clause so the
// reported total always matches the active filters.
func (r *taskRepository) List(filter TaskFilter) ([]domain.Task, int64, error) {
	q := filter.Query
	base := r.db.Model(&domain.Task{})

	if filter.UserID != nil {
		base = base.Where("user_id = ?", *filter.UserID)
	}
	if q.Status != "" {
		base = base.Where("status = ?", q.Status)
	}
	if q.Priority != "" {
		base = base.Where("priority = ?", q.Priority)
	}
	if q.Search != "" {
		// Case-insensitive title search. LOWER(...) LIKE keeps this portable
		// across PostgreSQL and SQLite.
		base = base.Where(`LOWER(title) LIKE ? ESCAPE '\'`, "%"+normalizeSearch(q.Search)+"%")
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var tasks []domain.Task
	err := base.
		Order(sortClause(q.Sort)).
		Offset(q.Offset()).
		Limit(q.Limit).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	return tasks, total, nil
}

// sortClause maps the validated sort token to a safe SQL ORDER BY expression.
// Because every branch is a hard-coded literal, no user input ever reaches the
// SQL string, eliminating injection risk on the ordering clause.
func sortClause(sort string) string {
	priorityRank := "CASE priority WHEN 'HIGH' THEN 3 WHEN 'MEDIUM' THEN 2 WHEN 'LOW' THEN 1 ELSE 0 END"
	switch sort {
	case "due_date_asc":
		return "due_date IS NULL, due_date ASC"
	case "due_date_desc":
		return "due_date IS NULL, due_date DESC"
	case "priority_asc":
		return priorityRank + " ASC"
	case "priority_desc":
		return priorityRank + " DESC"
	case "created_asc":
		return "created_at ASC"
	case "title_asc":
		return "LOWER(title) ASC"
	case "title_desc":
		return "LOWER(title) DESC"
	case "created_desc":
		fallthrough
	default:
		return "created_at DESC"
	}
}

// normalizeSearch lowercases the search term and escapes LIKE wildcards so a
// user searching for a literal "%" does not match everything.
func normalizeSearch(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		switch r {
		case '%', '_', '\\':
			out = append(out, '\\', r)
		default:
			out = append(out, toLowerRune(r))
		}
	}
	return string(out)
}

func toLowerRune(r rune) rune {
	if r >= 'A' && r <= 'Z' {
		return r + ('a' - 'A')
	}
	return r
}
