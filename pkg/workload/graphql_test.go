package workload_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/apexsudo/common/pkg/workload"
	"github.com/hellofresh/health-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/vektah/gqlparser/v2/ast"
)

type stubSchema struct{}

func (s *stubSchema) Schema() *ast.Schema { return &ast.Schema{} }
func (s *stubSchema) Complexity(_ context.Context, _, _ string, _ int, _ map[string]any) (int, bool) {
	return 0, false
}
func (s *stubSchema) Exec(_ context.Context) graphql.ResponseHandler {
	return func(_ context.Context) *graphql.Response { return &graphql.Response{} }
}

func TestWithGraphql_panicsOnNilSchema(t *testing.T) {
	t.Parallel()
	srv := newTestServer(t)

	assert.Panics(t, func() { srv.WithGraphql(nil) })
}

func TestWithGraphql_registersGraphQLRoute(t *testing.T) {
	t.Parallel()
	mux := http.NewServeMux()
	h, _ := health.New()
	srv := workload.New(mux, workload.Options{Health: h})
	srv.WithGraphql(&stubSchema{})

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/graphql", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	assert.NotEqual(t, http.StatusNotFound, rr.Code, "/graphql should be registered")
}
