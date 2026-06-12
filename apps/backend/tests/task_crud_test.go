package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskCRUDLifecycle(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("Owner", "owner@example.com", "supersecret123")

	// Create.
	rec := env.request(http.MethodPost, "/api/tasks", token, gin.H{
		"title": "Write report", "description": "Q3 numbers", "priority": "HIGH",
	})
	require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	var created struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}
	decode(t, rec, &created)
	assert.Equal(t, "Write report", created.Title)
	assert.Equal(t, "TODO", created.Status)
	assert.Equal(t, "HIGH", created.Priority)

	// Read.
	rec = env.request(http.MethodGet, "/api/tasks/"+created.ID, token, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Update status -> COMPLETED.
	rec = env.request(http.MethodPut, "/api/tasks/"+created.ID, token, gin.H{
		"status": "COMPLETED",
	})
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	var updated struct {
		Status string `json:"status"`
	}
	decode(t, rec, &updated)
	assert.Equal(t, "COMPLETED", updated.Status)

	// Delete.
	rec = env.request(http.MethodDelete, "/api/tasks/"+created.ID, token, nil)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Confirm gone.
	rec = env.request(http.MethodGet, "/api/tasks/"+created.ID, token, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestCreateTaskValidation(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("Owner", "owner@example.com", "supersecret123")

	// Missing title.
	rec := env.request(http.MethodPost, "/api/tasks", token, gin.H{"description": "no title"})
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

	// Invalid status enum.
	rec = env.request(http.MethodPost, "/api/tasks", token, gin.H{
		"title": "ok", "status": "NOPE",
	})
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestTaskOwnershipIsolation(t *testing.T) {
	env := newTestEnv(t)
	aliceToken, _ := env.signup("Alice", "alice@example.com", "supersecret123")
	bobToken, _ := env.signup("Bob", "bob@example.com", "supersecret123")

	// Alice creates a task.
	rec := env.request(http.MethodPost, "/api/tasks", aliceToken, gin.H{"title": "Alice secret"})
	require.Equal(t, http.StatusCreated, rec.Code)
	var task struct {
		ID string `json:"id"`
	}
	decode(t, rec, &task)

	// Bob cannot read it (404, not 403 — existence is not disclosed).
	rec = env.request(http.MethodGet, "/api/tasks/"+task.ID, bobToken, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Bob cannot update it.
	rec = env.request(http.MethodPut, "/api/tasks/"+task.ID, bobToken, gin.H{"title": "hijacked"})
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Bob cannot delete it.
	rec = env.request(http.MethodDelete, "/api/tasks/"+task.ID, bobToken, nil)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Bob's task list is empty (isolation).
	rec = env.request(http.MethodGet, "/api/tasks", bobToken, nil)
	assert.Equal(t, http.StatusOK, rec.Code)
	var bobTasks []map[string]interface{}
	decode(t, rec, &bobTasks)
	assert.Empty(t, bobTasks)
}
