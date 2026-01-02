package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pouyatavakoli/CodeStreaks-web/config"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/handler"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/infrastructure/database"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/infrastructure/scheduler"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/repository"
	"github.com/pouyatavakoli/CodeStreaks-web/internal/service"
	"github.com/pouyatavakoli/CodeStreaks-web/pkg/codeforces"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Println("Configuration loaded")

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	submissionRepo := repository.NewSubmissionRepository(db.DB)

	// Initialize Codeforces client
	cfClient := codeforces.NewClient(cfg.Codeforces.BaseURL)

	// Initialize services
	userService := service.NewUserService(userRepo, submissionRepo, cfClient)
	syncService := service.NewSyncService(
		userRepo,
		submissionRepo,
		cfClient,
		cfg.Codeforces.WorkerPoolSize,
	)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	healthHandler := handler.NewHealthHandler(db)

	// Setup router
	router := handler.NewRouter(userHandler, healthHandler)
	engine := router.Setup()

	// Initialize and start scheduler
	sched := scheduler.NewScheduler(syncService, cfg.Codeforces.UpdateInterval)
	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}
	defer sched.Stop()

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited successfully")
}
