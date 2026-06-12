package server

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskmanager/internal/config"
	"taskmanager/internal/handler"
	"taskmanager/internal/repository"
	"taskmanager/internal/service"
	"taskmanager/internal/sse"
	"taskmanager/internal/storage"
	"taskmanager/internal/utils"
)

// Container holds the fully wired application dependencies. Building it in one
// place makes the dependency graph explicit and lets tests construct the same
// graph against an in-memory database.
type Container struct {
	Config *config.Config
	DB     *gorm.DB

	JWT    *utils.JWTManager
	Broker *sse.Broker

	AuthHandler       *handler.AuthHandler
	TaskHandler       *handler.TaskHandler
	AttachmentHandler *handler.AttachmentHandler
	AdminHandler      *handler.AdminHandler
	SSEHandler        *handler.SSEHandler
}

// NewContainer constructs every repository, service, and handler from the given
// configuration and database connection.
func NewContainer(cfg *config.Config, db *gorm.DB) (*Container, error) {
	jwt := utils.NewJWTManager(cfg.JWTSecret, cfg.JWTExpiryHours)
	broker := sse.NewBroker()

	store, err := storage.NewLocalStorage(cfg.UploadDir)
	if err != nil {
		return nil, err
	}

	// Repositories.
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	activityRepo := repository.NewActivityRepository(db)
	attachmentRepo := repository.NewAttachmentRepository(db)

	// Services.
	authSvc := service.NewAuthService(userRepo, jwt)
	activitySvc := service.NewActivityService(activityRepo)
	taskSvc := service.NewTaskService(taskRepo, activitySvc, broker)
	attachmentSvc := service.NewAttachmentService(
		attachmentRepo, taskSvc, activitySvc, store, cfg.MaxUploadBytes, cfg.AllowedMime,
	)
	adminSvc := service.NewAdminService(userRepo)

	return &Container{
		Config:            cfg,
		DB:                db,
		JWT:               jwt,
		Broker:            broker,
		AuthHandler:       handler.NewAuthHandler(authSvc),
		TaskHandler:       handler.NewTaskHandler(taskSvc),
		AttachmentHandler: handler.NewAttachmentHandler(attachmentSvc),
		AdminHandler:      handler.NewAdminHandler(adminSvc, activitySvc),
		SSEHandler:        handler.NewSSEHandler(broker, jwt),
	}, nil
}

// Router builds the gin engine with the full middleware chain and route table.
func (ct *Container) Router() *gin.Engine {
	if ct.Config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.MaxMultipartMemory = ct.Config.MaxUploadBytes

	limiter := buildRateLimiter(ct.Config)

	registerMiddleware(r, ct.Config, limiter)
	registerRoutes(r, ct)
	return r
}
