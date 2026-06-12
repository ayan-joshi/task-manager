package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"taskmanager/internal/config"
	"taskmanager/internal/database"
	"taskmanager/internal/server"
)

// dbCounter gives each test its own uniquely-named in-memory database. SQLite's
// shared cache is keyed by the database name, so a unique name per test keeps
// the connection pool consistent while isolating tests from one another.
var dbCounter atomic.Int64

// testEnv bundles everything a test needs to drive the API end-to-end against an
// in-memory SQLite database — no Docker or external Postgres required.
type testEnv struct {
	t      *testing.T
	engine *gin.Engine
	db     *gorm.DB
	cfg    *config.Config
}

// newTestEnv builds a fully wired application backed by a fresh in-memory
// database. Each call is fully isolated.
func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	gin.SetMode(gin.TestMode)

	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared&_pragma=foreign_keys(1)", dbCounter.Add(1))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, database.Migrate(db))

	cfg := &config.Config{
		Port:           "8080",
		JWTSecret:      "test-secret-key-which-is-long-enough",
		JWTExpiryHours: 1,
		AllowedOrigins: []string{"*"},
		UploadDir:      t.TempDir(),
		MaxUploadBytes: 5 * 1024 * 1024,
		AllowedMime:    []string{"text/plain", "image/png"},
		RateLimitRPS:   1000,
		RateLimitBurst: 1000,
		Environment:    "test",
	}

	container, err := server.NewContainer(cfg, db)
	require.NoError(t, err)

	return &testEnv{t: t, engine: container.Router(), db: db, cfg: cfg}
}

// request performs an HTTP request against the in-memory engine and returns the
// recorder. A non-empty token is sent as a Bearer Authorization header.
func (e *testEnv) request(method, path, token string, body interface{}) *httptest.ResponseRecorder {
	e.t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(e.t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	e.engine.ServeHTTP(rec, req)
	return rec
}

// decode unmarshals the "data" field of the standard response envelope.
func decode(t *testing.T, rec *httptest.ResponseRecorder, target interface{}) {
	t.Helper()
	var env struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
		Meta    json.RawMessage `json:"meta"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &env))
	if target != nil && len(env.Data) > 0 {
		require.NoError(t, json.Unmarshal(env.Data, target))
	}
}

// signup registers a user and returns the issued token and user id.
func (e *testEnv) signup(name, email, password string) (token, userID string) {
	e.t.Helper()
	rec := e.request(http.MethodPost, "/api/auth/signup", "", gin.H{
		"name": name, "email": email, "password": password,
	})
	require.Equal(e.t, http.StatusCreated, rec.Code, rec.Body.String())
	var out struct {
		Token string `json:"token"`
		User  struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	decode(e.t, rec, &out)
	return out.Token, out.User.ID
}
