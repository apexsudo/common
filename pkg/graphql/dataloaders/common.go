// Package dataloaders provides generic dataloader primitives built on top of
// github.com/graph-gophers/dataloader/v7.
package dataloaders

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
)

const (
	// MaxBatchSize is the maximum number of keys collected into a single batch call.
	MaxBatchSize = 100
	// WaitTime is how long a loader waits to accumulate keys before dispatching a batch.
	WaitTime = 10 * time.Millisecond
)

// KeyFunc extracts the key from a value.
type KeyFunc[TKey comparable, TValue any] = func(TValue) TKey

// BatchFunc fetches a batch of values by their keys.
type BatchFunc[TKey comparable, TValue any] = func(ctx context.Context, keys []TKey) ([]TValue, error)

// Loader is the interface for a dataloader.
type Loader[TKey comparable, TValue any] interface {
	Load(ctx context.Context, key TKey) dataloader.Thunk[TValue]
}
