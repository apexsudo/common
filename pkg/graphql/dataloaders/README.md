# dataloaders

`dataloaders` provides generic, batched dataloader constructors built on top of [`graph-gophers/dataloader/v7`](https://github.com/graph-gophers/dataloader). Use these to solve the N+1 query problem in GraphQL resolvers.

All loaders are configured with sane defaults:

| Constant | Value | Purpose |
|----------|-------|---------|
| `MaxBatchSize` | 100 | Maximum keys collected before a batch is dispatched |
| `WaitTime` | 10 ms | How long to wait for additional keys before dispatching |

Cache is cleared after every batch (`WithClearCacheOnBatch`) so loaders are safe to reuse across requests when stored in context.

## Constructors

### `Single` — one-to-one relationships

Returns a `*dataloader.Loader[TKey, *TValue]`. Each key resolves to at most one value; missing keys resolve to `nil`.

```go
loader := dataloaders.Single(
    func(ctx context.Context, ids []int) ([]User, error) {
        return db.UsersByIDs(ctx, ids)
    },
    func(u User) int { return u.ID },
)

// Inside a resolver:
thunk := loader.Load(ctx, userID)
user, err := thunk()
```

### `Many` — one-to-many relationships

Returns a `*dataloader.Loader[TKey, []TValue]`. Each key resolves to a slice of values (empty slice for missing keys).

```go
loader := dataloaders.Many(
    func(ctx context.Context, userIDs []int) ([]Post, error) {
        return db.PostsByUserIDs(ctx, userIDs)
    },
    func(p Post) int { return p.UserID },
)

// Inside a resolver:
thunk := loader.Load(ctx, userID)
posts, err := thunk()
```

## Types

```go
// KeyFunc extracts the grouping key from a value.
type KeyFunc[TKey comparable, TValue any] = func(TValue) TKey

// BatchFunc fetches a batch of values by their keys.
type BatchFunc[TKey comparable, TValue any] = func(ctx context.Context, keys []TKey) ([]TValue, error)

// Loader is the common interface satisfied by both Single and Many loaders.
type Loader[TKey comparable, TValue any] interface {
    Load(ctx context.Context, key TKey) dataloader.Thunk[TValue]
}
```

## Usage pattern

Loaders are typically constructed once per request and stored in the request context so all resolvers in that request share the same batching window.

```go
type LoaderCtxKey struct{}

type Loaders struct {
    UserByID *dataloader.Loader[int, *User]
    PostsByUserID *dataloader.Loader[int, []Post]
}

func Middleware(db *DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            loaders := &Loaders{
                UserByID:      dataloaders.Single(db.UsersByIDs, func(u User) int { return u.ID }),
                PostsByUserID: dataloaders.Many(db.PostsByUserIDs, func(p Post) int { return p.UserID }),
            }
            ctx := context.WithValue(r.Context(), LoaderCtxKey{}, loaders)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```
