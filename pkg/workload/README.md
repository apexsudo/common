# workload

`workload` wires together multiple server workloads — HTTP, gRPC, GraphQL, or anything custom — under a single graceful-shutdown lifecycle powered by [`oklog/run`](https://github.com/oklog/run).

Every server gets three built-in HTTP routes for free:

| Route | Purpose |
|-------|---------|
| `GET /internal/healthz` | Health checks (hellofresh/health-go) |
| `GET /internal/metrics` | Prometheus metrics |
| `GET /internal/docs`    | Swagger UI |

## Quick start

```go
import "github.com/apexsudo/common/pkg/workload"

mux := http.NewServeMux()

srv := workload.New(mux, workload.Options{
    ServiceName: "my-service",
    Version:     "1.0.0",
    Health:      h, // *health.Health
})

srv.
    WithGracefulShutdown(ctx).
    WithGRPC(ctx, workload.GRPCOptions{EnableReflection: true}, grpcSvc).
    WithGraphql(executableSchema).
    Run(logger, nil)
}
```

## API

### `New(mux, opts) *Server`

Creates the server and registers the built-in internal routes on the provided mux.

### `(*Server) Run(logger, port) error`

Starts the HTTP server (default port **8000**) and all registered workloads concurrently. Blocks until every workload stops. The first workload to return an error triggers graceful shutdown of the rest.

### `(*Server) WithGracefulShutdown(ctx) *Server`

Listens for `SIGINT`/`SIGTERM`. When either arrives — or when `ctx` is cancelled — it initiates an orderly shutdown of the whole group.

### `(*Server) WithGRPC(ctx, opts, services...) *Server`

Adds a gRPC server (default port **50051**). Pass `GRPCOptions{EnableReflection: true}` to enable server reflection. Stopped gracefully alongside the rest of the group.

### `(*Server) WithGraphql(schema) *Server`

Registers a [gqlgen](https://gqlgen.com) handler at `/graphql` with POST, GET, and OPTIONS transports, an LRU query cache (1 000 entries), and introspection enabled. Panics if `schema` is nil.

### `(*Server) WithCustom(execute, interrupt) *Server`

Adds any arbitrary workload. `execute` runs the workload; `interrupt` is called when the group starts shutting down. Panics inside `execute` are caught and surfaced as errors.
