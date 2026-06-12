package service

import (
	"strings"

	"github.com/google/uuid"

	"taskmanager/internal/domain"
	"taskmanager/internal/dto"
	"taskmanager/internal/repository"
	"taskmanager/internal/sse"
	"taskmanager/internal/utils"
)

// TaskService implements task business logic: ownership enforcement, activity
// logging, and real-time event publication.
type TaskService struct {
	tasks    repository.TaskRepository
	activity *ActivityService
	broker   *sse.Broker
}

// NewTaskService constructs a TaskService.
func NewTaskService(tasks repository.TaskRepository, activity *ActivityService, broker *sse.Broker) *TaskService {
	return &TaskService{tasks: tasks, activity: activity, broker: broker}
}

// Create persists a new task owned by actorID.
func (s *TaskService) Create(actorID uuid.UUID, req dto.CreateTaskRequest) (*domain.Task, error) {
	task := &domain.Task{
		UserID:      actorID,
		Title:       strings.TrimSpace(req.Title),
		Description: req.Description,
		Status:      domain.Status(orDefault(req.Status, string(domain.StatusTodo))),
		Priority:    domain.Priority(orDefault(req.Priority, string(domain.PriorityMedium))),
		DueDate:     req.DueDate,
	}
	if err := s.tasks.Create(task); err != nil {
		return nil, err
	}

	s.activity.Record(actorID, &task.ID, domain.ActionTaskCreated, map[string]interface{}{
		"title": task.Title,
	})
	s.broker.Publish(actorID, dto.TaskEvent{Type: "task.created", TaskID: task.ID, Task: task})
	return task, nil
}

// GetByID returns a task if the actor is permitted to see it. Owners see their
// own tasks; admins see any task.
func (s *TaskService) GetByID(actorID uuid.UUID, role domain.Role, id uuid.UUID) (*domain.Task, error) {
	task, err := s.tasks.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, utils.ErrTaskNotFound
	}
	if !canAccess(actorID, role, task.UserID) {
		// Return "not found" rather than "forbidden" so the existence of other
		// users' task IDs is not disclosed.
		return nil, utils.ErrTaskNotFound
	}
	return task, nil
}

// List returns a filtered, paginated page of tasks. By default everyone —
// admins included — sees only their own tasks, so the personal board is never
// mixed with other users' work. An admin may opt into a global view with
// scope=all; for any non-admin that request silently falls back to own tasks.
func (s *TaskService) List(actorID uuid.UUID, role domain.Role, q dto.TaskQuery) ([]domain.Task, dto.PageMeta, error) {
	q.Normalize()
	filter := repository.TaskFilter{Query: q}

	allScope := role == domain.RoleAdmin && q.Scope == "all"
	if allScope {
		filter.IncludeOwner = true
	} else {
		filter.UserID = &actorID
	}

	tasks, total, err := s.tasks.List(filter)
	if err != nil {
		return nil, dto.PageMeta{}, err
	}
	return tasks, dto.NewPageMeta(total, q.Page, q.Limit), nil
}

// Update applies a partial update to a task the actor owns (or any task for an
// admin) and records the change.
func (s *TaskService) Update(actorID uuid.UUID, role domain.Role, id uuid.UUID, req dto.UpdateTaskRequest) (*domain.Task, error) {
	task, err := s.GetByID(actorID, role, id)
	if err != nil {
		return nil, err
	}

	statusChanged := false
	if req.Title != nil {
		task.Title = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Priority != nil {
		task.Priority = domain.Priority(*req.Priority)
	}
	if req.Status != nil && domain.Status(*req.Status) != task.Status {
		task.Status = domain.Status(*req.Status)
		statusChanged = true
	}
	if req.ClearDue {
		task.DueDate = nil
	} else if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.tasks.Update(task); err != nil {
		return nil, err
	}

	action := domain.ActionTaskUpdated
	if statusChanged {
		action = domain.ActionTaskStatusChanged
	}
	s.activity.Record(actorID, &task.ID, action, map[string]interface{}{
		"status": string(task.Status),
	})
	s.broker.Publish(task.UserID, dto.TaskEvent{Type: "task.updated", TaskID: task.ID, Task: task})
	return task, nil
}

// Delete removes a task the actor owns (or any task for an admin).
func (s *TaskService) Delete(actorID uuid.UUID, role domain.Role, id uuid.UUID) error {
	task, err := s.GetByID(actorID, role, id)
	if err != nil {
		return err
	}
	if err := s.tasks.Delete(task); err != nil {
		return err
	}
	s.activity.Record(actorID, &task.ID, domain.ActionTaskDeleted, map[string]interface{}{
		"title": task.Title,
	})
	s.broker.Publish(task.UserID, dto.TaskEvent{Type: "task.deleted", TaskID: task.ID})
	return nil
}

// Activity returns the activity log for a task the actor may access.
func (s *TaskService) Activity(actorID uuid.UUID, role domain.Role, id uuid.UUID, q dto.ActivityLogQuery) ([]domain.ActivityLog, dto.PageMeta, error) {
	if _, err := s.GetByID(actorID, role, id); err != nil {
		return nil, dto.PageMeta{}, err
	}
	return s.activity.ListByTask(id, q)
}

// canAccess reports whether an actor may operate on a resource owned by ownerID.
func canAccess(actorID uuid.UUID, role domain.Role, ownerID uuid.UUID) bool {
	return role == domain.RoleAdmin || actorID == ownerID
}

func orDefault(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
