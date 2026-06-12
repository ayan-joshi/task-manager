package tests

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedTasks creates a fixed set of tasks for filtering/sorting assertions.
func seedTasks(t *testing.T, env *testEnv, token string) {
	t.Helper()
	fixtures := []gin.H{
		{"title": "Write quarterly report", "status": "TODO", "priority": "HIGH"},
		{"title": "Review pull request", "status": "IN_PROGRESS", "priority": "MEDIUM"},
		{"title": "Plan sprint", "status": "TODO", "priority": "LOW"},
		{"title": "Deploy report service", "status": "COMPLETED", "priority": "HIGH"},
		{"title": "Update documentation", "status": "TODO", "priority": "MEDIUM"},
	}
	for _, f := range fixtures {
		rec := env.request(http.MethodPost, "/api/tasks", token, f)
		require.Equal(t, http.StatusCreated, rec.Code, rec.Body.String())
	}
}

type listResponse struct {
	data []struct {
		Title    string `json:"title"`
		Status   string `json:"status"`
		Priority string `json:"priority"`
	}
	meta struct {
		Total      int64 `json:"total"`
		Page       int   `json:"page"`
		Limit      int   `json:"limit"`
		TotalPages int   `json:"totalPages"`
	}
}

func (env *testEnv) list(t *testing.T, token string, query url.Values) listResponse {
	t.Helper()
	rec := env.request(http.MethodGet, "/api/tasks?"+query.Encode(), token, nil)
	require.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	var raw struct {
		Data json.RawMessage `json:"data"`
		Meta json.RawMessage `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &raw))

	var out listResponse
	require.NoError(t, json.Unmarshal(raw.Data, &out.data))
	require.NoError(t, json.Unmarshal(raw.Meta, &out.meta))
	return out
}

func TestFilterByStatus(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("User One", "u@example.com", "supersecret123")
	seedTasks(t, env, token)

	q := url.Values{}
	q.Set("status", "TODO")
	res := env.list(t, token, q)

	assert.Equal(t, int64(3), res.meta.Total)
	for _, task := range res.data {
		assert.Equal(t, "TODO", task.Status)
	}
}

func TestSearchAndFilterCombined(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("User One", "u@example.com", "supersecret123")
	seedTasks(t, env, token)

	// search="report" matches "Write quarterly report" (TODO) and "Deploy
	// report service" (COMPLETED). Adding status=TODO must narrow to one.
	q := url.Values{}
	q.Set("search", "report")
	q.Set("status", "TODO")
	res := env.list(t, token, q)

	require.Equal(t, int64(1), res.meta.Total)
	assert.Equal(t, "Write quarterly report", res.data[0].Title)
}

func TestSortByPriorityDescending(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("User One", "u@example.com", "supersecret123")
	seedTasks(t, env, token)

	q := url.Values{}
	q.Set("sort", "priority_desc")
	q.Set("limit", "100")
	res := env.list(t, token, q)

	require.NotEmpty(t, res.data)
	// First result must be a HIGH priority task.
	assert.Equal(t, "HIGH", res.data[0].Priority)
}

func TestPaginationMetadata(t *testing.T) {
	env := newTestEnv(t)
	token, _ := env.signup("User One", "u@example.com", "supersecret123")
	seedTasks(t, env, token)

	q := url.Values{}
	q.Set("page", "1")
	q.Set("limit", "2")
	res := env.list(t, token, q)

	assert.Equal(t, int64(5), res.meta.Total)
	assert.Equal(t, 2, res.meta.Limit)
	assert.Equal(t, 3, res.meta.TotalPages) // ceil(5/2)
	assert.Len(t, res.data, 2)
}
