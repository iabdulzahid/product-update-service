package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iabdulzahid/product-update-service/internal/handler"
	"github.com/iabdulzahid/product-update-service/internal/repository"
	"github.com/iabdulzahid/product-update-service/internal/service"
	"github.com/iabdulzahid/product-update-service/pkg"
	"github.com/iabdulzahid/product-update-service/pkg/queue"
)

func main() {
	// Load configuration
	cfg, err := pkg.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("cfg: %v\n", cfg)

	// Initialize store and queue
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(cfg.QueueSize)

	// Create context for workers (do NOT cancel immediately)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start worker pool
	ps := service.NewProductService(store, eq, cfg.Workers)
	ps.StartWorkers(ctx)
	log.Printf("Started %d workers", cfg.Workers)

	// HTTP handlers and router
	h := handler.NewProductHandler(store, eq)
	router := NewRouter(h)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("Server running on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	gracefulShutdown(srv, eq, cancel, 5*time.Second)
}

func gracefulShutdown(srv *http.Server, eq *queue.EventQueue, cancel context.CancelFunc, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Stop workers
	cancel()

	// Close the queue safely
	eq.Close()

	// Shutdown HTTP server
	ctx, shutdownCancel := context.WithTimeout(context.Background(), timeout)
	defer shutdownCancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
