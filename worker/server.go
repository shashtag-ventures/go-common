package worker

import (
	"log/slog"

	"github.com/hibiken/asynq"
)

// Server wraps asynq.Server to manage background task processing.
type Server struct {
	server *asynq.Server
	mux    *asynq.ServeMux
	logger *slog.Logger
}

func NewServer(redisAddr string, logger *slog.Logger) *Server {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			LogLevel: asynq.InfoLevel,
		},
	)

	return &Server{
		server: srv,
		mux:    asynq.NewServeMux(),
		logger: logger,
	}
}

// RegisterHandlers registers task handlers with the server mux.
func (s *Server) RegisterHandlers(handlers ...map[string]asynq.HandlerFunc) {
	for _, mapping := range handlers {
		for pattern, handler := range mapping {
			s.mux.HandleFunc(pattern, handler)
		}
	}
}

// Run starts the worker server.
func (s *Server) Run() error {
	s.logger.Info("Starting Asynq worker server...")
	return s.server.Run(s.mux)
}

// Shutdown stops the server gracefully.
func (s *Server) Shutdown() {
	s.logger.Info("Stopping Asynq worker server...")
	s.server.Shutdown()
}
