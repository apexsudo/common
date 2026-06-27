package workload

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/vektah/gqlparser/v2/ast"
)

// WithGraphql registers a gqlgen handler at POST|GET /graphql with query
// caching (LRU 1000 entries) and introspection enabled.  Panics if schema is nil.
func (s *Server) WithGraphql(schema graphql.ExecutableSchema) *Server {
	if schema == nil {
		panic("schema cannot be nil")
	}
	srv := handler.New(schema)
	srv.AroundFields(fieldMiddleware)
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GRAPHQL{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})

	s.mux.Handle("/graphql", srv)

	return s
}

func fieldMiddleware(ctx context.Context, next graphql.Resolver) (any, error) {
	return next(ctx)
}
