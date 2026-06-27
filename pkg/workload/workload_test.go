package workload_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hellofresh/health-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apexsudo/common/pkg/workload"
)

func newTestServer(t *testing.T) *workload.Server {
	t.Helper()
	h, err := health.New()
	require.NoError(t, err)

	return workload.New(http.NewServeMux(), workload.Options{Health: h})
}

func freePort(t *testing.T) *uint {
	t.Helper()
	lc := net.ListenConfig{}
	listener, err := lc.Listen(t.Context(), "tcp", ":0")
	require.NoError(t, err)
	require.NoError(t, listener.Close())
	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	require.True(t, ok, "unexpected address type: %T", listener.Addr())

	return new(uint(tcpAddr.Port))
}

func TestNew_registersHealthzRoute(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	h, _ := health.New()
	workload.New(mux, workload.Options{Health: h})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/internal/healthz", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusNotFound, rr.Code, "/internal/healthz should be registered")
}

func TestNew_registersMetricsRoute(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	h, _ := health.New()
	workload.New(mux, workload.Options{Health: h})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/internal/metrics", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusNotFound, rr.Code, "/internal/metrics should be registered")
}

func TestNew_registersDocsRoute(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	h, _ := health.New()
	workload.New(mux, workload.Options{Health: h})

	// POST triggers 405 (not 404) when the GET-prefixed route exists in Go 1.22+ routing.
	req := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/internal/docs", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "/internal/docs should be registered as GET route")
}
