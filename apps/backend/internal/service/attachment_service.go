package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"

	"taskmanager/internal/domain"
	"taskmanager/internal/repository"
	"taskmanager/internal/storage"
	"taskmanager/internal/utils"
)

// AttachmentService handles upload, retrieval, and deletion of task attachments,
// enforcing the same ownership rules as tasks plus file type/size validation.
type AttachmentService struct {
	repo        repository.AttachmentRepository
	tasks       *TaskService
	activity    *ActivityService
	store       storage.Storage
	maxBytes    int64
	allowedMime map[string]struct{}
}

// NewAttachmentService constructs an AttachmentService.
func NewAttachmentService(
	repo repository.AttachmentRepository,
	tasks *TaskService,
	activity *ActivityService,
	store storage.Storage,
	maxBytes int64,
	allowedMime []string,
) *AttachmentService {
	set := make(map[string]struct{}, len(allowedMime))
	for _, m := range allowedMime {
		set[m] = struct{}{}
	}
	return &AttachmentService{
		repo:        repo,
		tasks:       tasks,
		activity:    activity,
		store:       store,
		maxBytes:    maxBytes,
		allowedMime: set,
	}
}

// Upload validates and stores a file, associating it with a task the actor owns.
func (s *AttachmentService) Upload(
	actorID uuid.UUID,
	role domain.Role,
	taskID uuid.UUID,
	header *multipart.FileHeader,
	file multipart.File,
) (*domain.Attachment, error) {
	// Ownership: actor must be able to access the parent task.
	task, err := s.tasks.GetByID(actorID, role, taskID)
	if err != nil {
		return nil, err
	}

	if header.Size > s.maxBytes {
		return nil, utils.ErrFileTooLarge
	}

	contentType := header.Header.Get("Content-Type")
	if _, ok := s.allowedMime[contentType]; !ok {
		return nil, utils.ErrUnsupportedType
	}

	storedName := fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(header.Filename))
	if err := s.store.Save(storedName, file); err != nil {
		return nil, err
	}

	att := &domain.Attachment{
		TaskID:      task.ID,
		UserID:      task.UserID,
		FileName:    filepath.Base(header.Filename),
		StoredName:  storedName,
		ContentType: contentType,
		SizeBytes:   header.Size,
	}
	if err := s.repo.Create(att); err != nil {
		// Roll back the stored bytes if the DB row could not be created.
		_ = s.store.Delete(storedName)
		return nil, err
	}

	s.activity.Record(actorID, &task.ID, domain.ActionAttachmentAdded, map[string]interface{}{
		"fileName": att.FileName,
	})
	return att, nil
}

// Open returns an attachment's metadata and a reader for its bytes, after
// verifying the actor may access the parent task.
func (s *AttachmentService) Open(actorID uuid.UUID, role domain.Role, id uuid.UUID) (*domain.Attachment, io.ReadCloser, error) {
	att, err := s.repo.FindByID(id)
	if err != nil {
		return nil, nil, err
	}
	if att == nil {
		return nil, nil, utils.ErrAttachNotFound
	}
	if _, err := s.tasks.GetByID(actorID, role, att.TaskID); err != nil {
		return nil, nil, utils.ErrAttachNotFound
	}
	reader, err := s.store.Open(att.StoredName)
	if err != nil {
		return nil, nil, utils.ErrAttachNotFound
	}
	return att, reader, nil
}

// Delete removes an attachment and its stored bytes.
func (s *AttachmentService) Delete(actorID uuid.UUID, role domain.Role, id uuid.UUID) error {
	att, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if att == nil {
		return utils.ErrAttachNotFound
	}
	if _, err := s.tasks.GetByID(actorID, role, att.TaskID); err != nil {
		return utils.ErrAttachNotFound
	}
	if err := s.repo.Delete(att); err != nil {
		return err
	}
	_ = s.store.Delete(att.StoredName)
	s.activity.Record(actorID, &att.TaskID, domain.ActionAttachmentDeleted, map[string]interface{}{
		"fileName": att.FileName,
	})
	return nil
}
