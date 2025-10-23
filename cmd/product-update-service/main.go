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

	"github.com/iabdulzahid/product-update-service/internal/config"
	"github.com/iabdulzahid/product-update-service/internal/handler"
	"github.com/iabdulzahid/product-update-service/internal/repository"
	"github.com/iabdulzahid/product-update-service/internal/service"
	"github.com/iabdulzahid/product-update-service/pkg/queue"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("cfg: %v\n", cfg)
	// Initialize store and queue
	store := repository.NewProductStore()
	eq := queue.NewEventQueue(cfg.QueueSize)

	// Start worker pool
	ps := service.NewProductService(store, eq, cfg.Workers)
	// ps.StartWorkers()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ps := service.NewProductService(store, eq, cfg.Workers)
	ps.StartWorkers(ctx)

	// on graceful shutdown:
	cancel()   // signal workers to stop (optional)
	eq.Close() // implement Close() to close the channel safely
	// wait for a short period for workers to finish

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
	gracefulShutdown(srv, 5*time.Second)
}

// gracefulShutdown handles SIGINT/SIGTERM and shuts down the server gracefully
func gracefulShutdown(srv *http.Server, timeout time.Duration) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
