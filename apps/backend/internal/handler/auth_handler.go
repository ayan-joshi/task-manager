package handler

import (
	"github.com/gin-gonic/gin"

	"taskmanager/internal/dto"
	"taskmanager/internal/middleware"
	"taskmanager/internal/service"
	"taskmanager/internal/utils"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler constructs an AuthHandler.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Signup handles POST /api/auth/signup.
func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if !bindJSON(c, &req) {
		return
	}
	res, err := h.auth.Signup(req)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.Created(c, res)
}

// Login handles POST /api/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if !bindJSON(c, &req) {
		return
	}
	res, err := h.auth.Login(req)
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, res)
}

// Me handles GET /api/auth/me, returning the current authenticated user.
func (h *AuthHandler) Me(c *gin.Context) {
	user, err := h.auth.Me(middleware.UserID(c))
	if err != nil {
		utils.Fail(c, err)
		return
	}
	utils.OK(c, user)
}
