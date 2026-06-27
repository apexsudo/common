// Package workload provides a composable HTTP server with optional gRPC, GraphQL,
// and custom workloads that all share a single graceful-shutdown lifecycle managed
// by oklog/run.  Built-in routes are registered automatically:
//
//   - GET /internal/healthz  — health checks (hellofresh/health-go)
//   - GET /internal/metrics  — Prometheus metrics
//   - GET /internal/docs     — Swagger UI
//
// Typical usage:
//
//	srv := workload.New(mux, workload.Options{...})
//	srv.WithGracefulShutdown(ctx).
//	    WithGRPC(ctx, workload.GRPCOptions{}, myGRPCSvc).
//	    WithGraphql(mySchema)
//	if err := srv.Run(logger, nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
//	    log.Fatal(err)
//	}
package workload

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/samber/lo"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Options configures a Server.
type Options struct {
	ServiceName string
	Version     string
	Health      *health.Health
}

// Server is the central workload coordinator.  Add workloads via With* methods,
// then call Run to start all of them under a shared lifecycle.
type Server struct {
	opts   Options
	mux    *http.ServeMux
	group  *run.Group
	logger *slog.Logger
}

// New creates a Server, registers the built-in internal routes on mux, and
// returns it ready for With* calls.
func New(mux *http.ServeMux, opts Options) *Server {
	mux.Handle("GET /internal/healthz", opts.Health.Handler())
	mux.Handle("GET /internal/metrics", promhttp.Handler())
	mux.Handle("GET /internal/docs", httpSwagger.Handler())

	return &Server{
		opts:   opts,
		mux:    mux,
		group:  &run.Group{},
		logger: slog.Default(),
	}
}

// Run starts the HTTP server on port (default 8000) together with all
// registered workloads.  It blocks until every workload has stopped.
// The first workload to return an error stops the rest via graceful shutdown.
func (s *Server) Run(logger *slog.Logger, port *uint) error {
	s.logger = logger

	httpAddr := lo.FromPtrOr(port, 8000)
	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", httpAddr),
		Handler:           s.mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	s.group.Add(
		func() error {
			logger.Info("HTTP listening", "addr", httpAddr)

			return httpServer.ListenAndServe()
		},
		func(interruptReason error) {
			shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			s.logger.Info("Shutting down HTTP server", "reason", interruptReason)
			err := httpServer.Shutdown(shutCtx)
			if err != nil {
				s.logger.Warn("Failed to shutdown HTTP server", "err", err)
			}
		},
	)

	logger.Info("starting", "service", s.opts.ServiceName, "version", s.opts.Version)

	if err := s.group.Run(); err != nil {
		return fmt.Errorf("workload: %w", err)
	}

	return nil
}
