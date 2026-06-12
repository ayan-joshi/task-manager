// Command api is the entrypoint for the Task Manager backend HTTP server.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskmanager/internal/config"
	"taskmanager/internal/database"
	"taskmanager/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	db, err := database.Connect(cfg.DatabaseURL, cfg.IsProduction())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	if err := database.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	if err := server.SeedAdmin(db, cfg); err != nil {
		log.Fatalf("admin seed failed: %v", err)
	}

	container, err := server.NewContainer(cfg, db)
	if err != nil {
		log.Fatalf("dependency wiring failed: %v", err)
	}

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           container.Router(),
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Run the server in a goroutine so the main goroutine can wait for signals.
	go func() {
		log.Printf("server listening on :%s (env=%s)", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM gives in-flight requests time to finish.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
	log.Println("server stopped")
}
