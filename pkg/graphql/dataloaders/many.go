package dataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
)

// Many creates a batched dataloader for one-to-many relationships.
// batchFunc is called once per batch with all requested keys and must return
// all matching values in any order; keyFunc extracts the key from each value
// so results can be grouped back to their callers.
// A nil slice is returned for any key that has no matching values.
func Many[TKey comparable, TValue any](
	batchFunc BatchFunc[TKey, TValue],
	keyFunc KeyFunc[TKey, TValue],
) *dataloader.Loader[TKey, []TValue] {
	return dataloader.NewBatchedLoader[TKey, []TValue](
		func(ctx context.Context, keys []TKey) []*dataloader.Result[[]TValue] {
			values, err := batchFunc(ctx, keys)
			if err != nil {
				return lo.Map(keys, func(item TKey, index int) *dataloader.Result[[]TValue] {
					return &dataloader.Result[[]TValue]{Error: err}
				})
			}
			valuesMap := lo.GroupBy(values, func(item TValue) TKey {
				return keyFunc(item)
			})

			return lo.Map(keys, func(item TKey, index int) *dataloader.Result[[]TValue] {
				value, ok := valuesMap[item]
				if ok {
					return &dataloader.Result[[]TValue]{
						Data: value,
					}
				}

				return &dataloader.Result[[]TValue]{
					Data: nil,
				}
			})
		},
		dataloader.WithBatchCapacity[TKey, []TValue](MaxBatchSize),
		dataloader.WithWait[TKey, []TValue](WaitTime),
		dataloader.WithClearCacheOnBatch[TKey, []TValue](),
	)
}
