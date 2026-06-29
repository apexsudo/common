package workload

import (
	"context"
	"html/template"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

var uiIndexPage = template.Must(template.New("ui").Parse(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>GraphQL Explorers</title>
  <style>
    *, *::before, *::after { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
      background: #0f0f11;
      color: #e8e8ea;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      gap: 12px;
    }
    h1 { font-size: 1.1rem; font-weight: 500; color: #888; margin: 0 0 20px; letter-spacing: .05em; text-transform: uppercase; }
    .cards { display: flex; gap: 16px; flex-wrap: wrap; justify-content: center; }
    a.card {
      display: flex;
      flex-direction: column;
      gap: 6px;
      width: 200px;
      padding: 24px 20px;
      background: #1a1a1f;
      border: 1px solid #2a2a30;
      border-radius: 10px;
      text-decoration: none;
      color: inherit;
      transition: border-color .15s, background .15s;
    }
    a.card:hover { border-color: #555; background: #1f1f26; }
    .card-name { font-size: 1rem; font-weight: 600; }
    .card-desc { font-size: .8rem; color: #666; }
  </style>
</head>
<body>
  <h1>Choose a GraphQL Explorer</h1>
  <div class="cards">
    <a class="card" href="/ui/graphiql">
      <span class="card-name">GraphiQL</span>
      <span class="card-desc">The classic in-browser IDE by the GraphQL Foundation</span>
    </a>
    <a class="card" href="/ui/apollo">
      <span class="card-name">Apollo Sandbox</span>
      <span class="card-desc">Apollo Studio embedded explorer</span>
    </a>
    <a class="card" href="/ui/altair">
      <span class="card-name">Altair</span>
      <span class="card-desc">Feature-rich GraphQL client with subscriptions</span>
    </a>
  </div>
</body>
</html>`))

// WithGraphql registers a gqlgen handler at POST|GET /graphql with query
// caching (LRU 1000 entries) and introspection enabled. Three playground UIs
// are mounted under /ui and a picker page is served at /ui. Panics if schema is nil.
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

	s.mux.Handle("/ui/graphiql", playground.Handler("GraphiQL", "/graphql"))
	s.mux.Handle("/ui/apollo", playground.ApolloSandboxHandler("Apollo Sandbox", "/graphql"))
	s.mux.Handle("/ui/altair", playground.AltairHandler("Altair", "/graphql", nil))

	s.mux.HandleFunc("/ui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		if err := uiIndexPage.Execute(w, nil); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	})
	s.mux.HandleFunc("/ui/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/ui", http.StatusMovedPermanently)
	})

	return s
}

func fieldMiddleware(ctx context.Context, next graphql.Resolver) (any, error) {
	return next(ctx)
}
