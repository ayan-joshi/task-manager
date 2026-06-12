package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"taskmanager/internal/domain"
	"taskmanager/internal/utils"
)

// createAdmin inserts an ADMIN user directly and returns a valid token for it.
func (e *testEnv) createAdmin(email, password string) string {
	e.t.Helper()
	hash, err := utils.HashPassword(password)
	require.NoError(e.t, err)
	admin := &domain.User{Name: "Admin", Email: email, Password: hash, Role: domain.RoleAdmin}
	require.NoError(e.t, e.db.Create(admin).Error)

	rec := e.request(http.MethodPost, "/api/auth/login", "", gin.H{"email": email, "password": password})
	require.Equal(e.t, http.StatusOK, rec.Code, rec.Body.String())
	var out struct {
		Token string `json:"token"`
	}
	decode(e.t, rec, &out)
	return out.Token
}

func TestAdminTaskScope(t *testing.T) {
	env := newTestEnv(t)

	// Two regular users each own a task.
	aliceToken, _ := env.signup("Alice", "alice@example.com", "supersecret123")
	bobToken, _ := env.signup("Bob", "bob@example.com", "supersecret123")
	require.Equal(t, http.StatusCreated,
		env.request(http.MethodPost, "/api/tasks", aliceToken, gin.H{"title": "Alice task"}).Code)
	require.Equal(t, http.StatusCreated,
		env.request(http.MethodPost, "/api/tasks", bobToken, gin.H{"title": "Bob task"}).Code)

	adminToken := env.createAdmin("admin@example.com", "supersecret123")

	// Admin's personal board (default scope) is empty — admin owns no tasks.
	rec := env.request(http.MethodGet, "/api/tasks", adminToken, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	var ownTasks []map[string]interface{}
	decode(t, rec, &ownTasks)
	assert.Empty(t, ownTasks, "admin default list must be own-scoped")

	// With scope=all, the admin sees every user's tasks, each with its owner.
	rec = env.request(http.MethodGet, "/api/tasks?scope=all", adminToken, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	var allTasks []struct {
		Title string `json:"title"`
		Owner *struct {
			Email string `json:"email"`
		} `json:"owner"`
	}
	decode(t, rec, &allTasks)
	require.Len(t, allTasks, 2)
	assert.NotNil(t, allTasks[0].Owner, "admin all-scope list should include task owner")

	// A non-admin requesting scope=all is silently confined to their own tasks.
	rec = env.request(http.MethodGet, "/api/tasks?scope=all", aliceToken, nil)
	require.Equal(t, http.StatusOK, rec.Code)
	var aliceAll []map[string]interface{}
	decode(t, rec, &aliceAll)
	assert.Len(t, aliceAll, 1, "scope=all must not escalate a non-admin")
}
