package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignupAndLogin(t *testing.T) {
	env := newTestEnv(t)

	// Signup succeeds and returns a token.
	token, userID := env.signup("Ada Lovelace", "ada@example.com", "supersecret123")
	require.NotEmpty(t, token)
	require.NotEmpty(t, userID)

	// Duplicate email is rejected with 409.
	rec := env.request(http.MethodPost, "/api/auth/signup", "", gin.H{
		"name": "Ada Again", "email": "ada@example.com", "password": "supersecret123",
	})
	assert.Equal(t, http.StatusConflict, rec.Code)

	// Login with correct credentials succeeds.
	rec = env.request(http.MethodPost, "/api/auth/login", "", gin.H{
		"email": "ada@example.com", "password": "supersecret123",
	})
	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())

	// Login with wrong password fails with 401.
	rec = env.request(http.MethodPost, "/api/auth/login", "", gin.H{
		"email": "ada@example.com", "password": "wrongpassword",
	})
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestSignupValidation(t *testing.T) {
	env := newTestEnv(t)

	// Invalid email and short password produce a 422 validation error.
	rec := env.request(http.MethodPost, "/api/auth/signup", "", gin.H{
		"name": "X", "email": "not-an-email", "password": "short",
	})
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestProtectedRouteRequiresToken(t *testing.T) {
	env := newTestEnv(t)

	// No token: 401.
	rec := env.request(http.MethodGet, "/api/tasks", "", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	// Garbage token: 401.
	rec = env.request(http.MethodGet, "/api/tasks", "garbage.token.value", nil)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
