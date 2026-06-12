package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"taskmanager/internal/domain"
)

// AttachmentRepository abstracts persistence for task attachments.
type AttachmentRepository interface {
	Create(att *domain.Attachment) error
	FindByID(id uuid.UUID) (*domain.Attachment, error)
	Delete(att *domain.Attachment) error
}

type attachmentRepository struct{ db *gorm.DB }

// NewAttachmentRepository constructs an AttachmentRepository backed by GORM.
func NewAttachmentRepository(db *gorm.DB) AttachmentRepository {
	return &attachmentRepository{db: db}
}

func (r *attachmentRepository) Create(att *domain.Attachment) error {
	return r.db.Create(att).Error
}

func (r *attachmentRepository) FindByID(id uuid.UUID) (*domain.Attachment, error) {
	var att domain.Attachment
	err := r.db.First(&att, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *attachmentRepository) Delete(att *domain.Attachment) error {
	return r.db.Delete(att).Error
}
