package service

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"

	"taskmanager/internal/domain"
	"taskmanager/internal/dto"
	"taskmanager/internal/repository"
)

// ActivityService records and retrieves audit log entries.
type ActivityService struct {
	repo repository.ActivityRepository
}

// NewActivityService constructs an ActivityService.
func NewActivityService(repo repository.ActivityRepository) *ActivityService {
	return &ActivityService{repo: repo}
}

// Record writes an activity entry. It is best-effort: a logging failure must
// never break the user-facing operation that triggered it, so the error is
// logged and swallowed.
func (s *ActivityService) Record(userID uuid.UUID, taskID *uuid.UUID, action domain.Action, metadata map[string]interface{}) {
	var meta string
	if metadata != nil {
		if b, err := json.Marshal(metadata); err == nil {
			meta = string(b)
		}
	}
	entry := &domain.ActivityLog{
		UserID:   userID,
		TaskID:   taskID,
		Action:   action,
		Metadata: meta,
	}
	if err := s.repo.Create(entry); err != nil {
		log.Printf("activity log write failed: %v", err)
	}
}

// List returns a page of activity log entries (admin use).
func (s *ActivityService) List(q dto.ActivityLogQuery) ([]domain.ActivityLog, dto.PageMeta, error) {
	q.Normalize()
	logs, total, err := s.repo.List(q.Offset(), q.Limit)
	if err != nil {
		return nil, dto.PageMeta{}, err
	}
	return logs, dto.NewPageMeta(total, q.Page, q.Limit), nil
}

// ListByTask returns a page of activity entries for a single task. The caller is
// responsible for authorizing access to that task.
func (s *ActivityService) ListByTask(taskID uuid.UUID, q dto.ActivityLogQuery) ([]domain.ActivityLog, dto.PageMeta, error) {
	q.Normalize()
	logs, total, err := s.repo.ListByTask(taskID, q.Offset(), q.Limit)
	if err != nil {
		return nil, dto.PageMeta{}, err
	}
	return logs, dto.NewPageMeta(total, q.Page, q.Limit), nil
}
