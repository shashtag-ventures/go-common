package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// shutdownWG is used to wait for all goroutines to complete before exiting.
var shutdownWG sync.WaitGroup

// Server holds the dependencies for an HTTP server.
type Server struct {
	*http.Server
	shutdownTimeout time.Duration
}

// New creates a new Server instance.
func New(addr string, handler http.Handler, shutdownTimeout time.Duration) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    ":" + addr,
			Handler: handler,
		},
		shutdownTimeout: shutdownTimeout,
	}
}

// Start runs the server and handles graceful shutdown.
func (srv *Server) Start() {
	slog.Info("Server is running on port " + srv.Addr)

	// Increment the WaitGroup counter for the server goroutine.
	AddGracefulShutdownGoroutine()
	go func() {
		defer DoneGracefulShutdownGoroutine() // Decrement the counter when the goroutine exits.
		// ListenAndServe blocks until an error occurs or the server is shut down.
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to listen and serve", "error", err)
		}
	}()

	// Graceful shutdown logic
	shutdown, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()

	slog.Info("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down server in time", "error", err)
	}
	slog.Info("Server shutdown gracefully")

	slog.Info("Waiting for all goroutines to be completed")
	shutdownWG.Wait()
	slog.Info("All goroutines are completed")
}

// AddGracefulShutdownGoroutine adds a goroutine to the WaitGroup.
func AddGracefulShutdownGoroutine() {
	shutdownWG.Add(1)
}

// DoneGracefulShutdownGoroutine marks a goroutine as complete.
func DoneGracefulShutdownGoroutine() {
	shutdownWG.Done()
}
