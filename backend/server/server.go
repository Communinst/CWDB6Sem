package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}

func (s *Server) Shutdown(c context.Context) error {
	return s.httpServer.Shutdown(c)
}

func New(address string, handler http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         address,
			Handler:      handler,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		},
	}
}

func (s *Server) Run() {
	slog.Info("Server is running at", "address", s.httpServer.Addr)

	shutdownChan := make(chan bool, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP Server error: %v", err)
		}

		log.Println("Stopped serving new connection.")
		shutdownChan <- true
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), s.httpServer.ReadTimeout)
	defer shutdownRelease()

	if err := s.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server shutdown error: %v", err)
	}

	<-shutdownChan
	log.Println("Graceful shutdown complete.")
}
