package server

import (
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestServer_StartAndShutdown(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := New("8081", handler, 5*time.Second)

	go func() {
		// Give the server a moment to start
		time.Sleep(100 * time.Millisecond)
		// Send a signal to shut down
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	// This will block until shutdown is complete
	srv.Start()

	// If Start() returns, it means the server has shut down.
	// No explicit assertion is needed here, the test will pass if it doesn't hang.
}
