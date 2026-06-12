package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"taskmanager/internal/domain"
)

// ActivityRepository abstracts persistence for activity logs.
type ActivityRepository interface {
	Create(log *domain.ActivityLog) error
	List(offset, limit int) ([]domain.ActivityLog, int64, error)
	ListByTask(taskID uuid.UUID, offset, limit int) ([]domain.ActivityLog, int64, error)
}

type activityRepository struct{ db *gorm.DB }

// NewActivityRepository constructs an ActivityRepository backed by GORM.
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) Create(log *domain.ActivityLog) error {
	return r.db.Create(log).Error
}

func (r *activityRepository) List(offset, limit int) ([]domain.ActivityLog, int64, error) {
	var logs []domain.ActivityLog
	var total int64

	if err := r.db.Model(&domain.ActivityLog{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.
		Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

func (r *activityRepository) ListByTask(taskID uuid.UUID, offset, limit int) ([]domain.ActivityLog, int64, error) {
	var logs []domain.ActivityLog
	var total int64

	base := r.db.Model(&domain.ActivityLog{}).Where("task_id = ?", taskID)
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := base.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}
