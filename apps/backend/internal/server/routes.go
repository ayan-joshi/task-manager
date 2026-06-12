package server

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"taskmanager/internal/config"
	"taskmanager/internal/domain"
	"taskmanager/internal/handler"
	"taskmanager/internal/middleware"
)

func buildRateLimiter(cfg *config.Config) *middleware.RateLimiter {
	return middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
}

// registerMiddleware installs the global middleware chain in execution order:
// request ID → logging → panic recovery → security headers → CORS → rate limit.
func registerMiddleware(r *gin.Engine, cfg *config.Config, limiter *middleware.RateLimiter) {
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(limiter.Middleware())
}

// registerRoutes defines the full API surface.
func registerRoutes(r *gin.Engine, ct *Container) {
	api := r.Group("/api")

	api.GET("/health", handler.Health)
	// SSE authenticates via query-param token (EventSource cannot set headers).
	api.GET("/events", ct.SSEHandler.Stream)

	auth := api.Group("/auth")
	{
		auth.POST("/signup", ct.AuthHandler.Signup)
		auth.POST("/login", ct.AuthHandler.Login)
		auth.GET("/me", middleware.Auth(ct.JWT), ct.AuthHandler.Me)
	}

	// Everything below requires a valid JWT.
	protected := api.Group("")
	protected.Use(middleware.Auth(ct.JWT))
	{
		tasks := protected.Group("/tasks")
		{
			tasks.POST("", ct.TaskHandler.Create)
			tasks.GET("", ct.TaskHandler.List)
			tasks.GET("/:id", ct.TaskHandler.Get)
			tasks.PUT("/:id", ct.TaskHandler.Update)
			tasks.DELETE("/:id", ct.TaskHandler.Delete)
			tasks.GET("/:id/activity", ct.TaskHandler.Activity)
			tasks.POST("/:id/attachments", ct.AttachmentHandler.Upload)
		}

		attachments := protected.Group("/attachments")
		{
			attachments.GET("/:attachmentId", ct.AttachmentHandler.Download)
			attachments.DELETE("/:attachmentId", ct.AttachmentHandler.Delete)
		}

		admin := protected.Group("/admin")
		admin.Use(middleware.RequireRole(domain.RoleAdmin))
		{
			admin.GET("/users", ct.AdminHandler.ListUsers)
			admin.GET("/activity-logs", ct.AdminHandler.ListActivityLogs)
		}
	}
}
