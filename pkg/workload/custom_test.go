package workload_test

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithCustom_returnsServerForChaining(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)

	result := srv.WithCustom(func() error { return nil }, func(error) {})

	assert.Same(t, srv, result)
}

func TestWithCustom_recoversPanic(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	srv.WithCustom(func() error { panic("intentional panic") }, func(error) {})

	// The panic actor returns immediately; run.Group then shuts down the HTTP server.
	err := srv.Run(slog.Default(), freePort(t))

	require.Error(t, err)
	assert.ErrorContains(t, err, "panic")
}

func TestWithGracefulShutdown_stopsOnContextCancel(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel before Run so the graceful shutdown actor exits immediately
	srv.WithGracefulShutdown(ctx)

	done := make(chan error, 1)
	go func() { done <- srv.Run(slog.Default(), freePort(t)) }()

	select {
	case <-done:
		// server stopped as expected
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop after context cancellation")
	}
}
