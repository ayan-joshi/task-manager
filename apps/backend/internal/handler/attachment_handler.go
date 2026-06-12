package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"taskmanager/internal/middleware"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

// AttachmentHandler exposes file attachment endpoints.
type AttachmentHandler struct {
	attachments *service.AttachmentService
}

// NewAttachmentHandler constructs an AttachmentHandler.
func NewAttachmentHandler(attachments *service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{attachments: attachments}
}

// Upload handles POST /api/tasks/:id/attachments (multipart/form-data, field "file").
func (h *AttachmentHandler) Upload(c *gin.Context) {
	taskID, ok := parseUUIDParam(c, "id")
	if !ok {
		return
	}
	header, err := c.FormFile("file")
	if err != nil {
		utils.FailValidation(c, "missing 'file' form field")
		return
	}
	file, err := header.Open()
	if err != nil {
		utils.Fail(c, utils.ErrInternal)
		return
	}
	defer file.Close()

	att, err := h.attachments.Upload(
		middleware.UserID(c), middleware.CurrentRole(c), taskID, header, file,
	)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Created(c, att)
}

// Download handles GET /api/attachments/:attachmentId, streaming the file bytes.
func (h *AttachmentHandler) Download(c *gin.Context) {
	id, ok := parseUUIDParam(c, "attachmentId")
	if !ok {
		return
	}
	att, reader, err := h.attachments.Open(middleware.UserID(c), middleware.CurrentRole(c), id)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", att.ContentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", att.FileName))
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}

// Delete handles DELETE /api/attachments/:attachmentId.
func (h *AttachmentHandler) Delete(c *gin.Context) {
	id, ok := parseUUIDParam(c, "attachmentId")
	if !ok {
		return
	}
	if err := h.attachments.Delete(middleware.UserID(c), middleware.CurrentRole(c), id); err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, gin.H{"deleted": true, "id": id})
}
