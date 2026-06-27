package workload_test

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apexsudo/common/pkg/workload"
)

func TestWithGRPC_returnsServerForChaining(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	port := freePort(t)

	result := srv.WithGRPC(context.Background(), workload.GRPCOptions{Port: port})

	assert.Same(t, srv, result)
}

func TestWithGRPC_listensOnSpecifiedPort(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)
	grpcPort := freePort(t)

	srv.WithGRPC(t.Context(), workload.GRPCOptions{Port: grpcPort})
	srv.WithGracefulShutdown(t.Context())
	go srv.Run(slog.Default(), freePort(t)) //nolint:errcheck

	addr := fmt.Sprintf(":%d", *grpcPort)
	assert.Eventually(t, func() bool {
		conn, err := (&net.Dialer{}).DialContext(t.Context(), "tcp", addr)
		if err != nil {
			return false
		}
		require.NoError(t, conn.Close())

		return true
	}, 2*time.Second, 10*time.Millisecond, "gRPC server did not start listening on port %d", *grpcPort)
}
