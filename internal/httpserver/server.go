package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/guverz/pr-reviewer-service/internal/config"
)

type Server struct {
	httpServer      *http.Server
	shutdownTimeout time.Duration
}

func New(cfg *config.Config, router http.Handler) (*Server, error) {
	server := &http.Server{
		Addr:              cfg.HTTP.Addr,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTP.ReadHeader,
	}

	return &Server{
		httpServer:      server,
		shutdownTimeout: cfg.HTTP.ShutdownTimeout,
	}, nil
}

func (s *Server) Serve(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("start http server: %w", err)
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil
	case err := <-errCh:
		return err
	}
}
