package dataloaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/samber/lo"
)

// Single creates a batched dataloader for one-to-one relationships.
// batchFunc is called once per batch with all requested keys and must return
// the matching values in any order; keyFunc extracts the key from each value
// so each result can be routed back to its caller.
// A nil pointer is returned for any key that has no matching value.
func Single[TKey comparable, TValue any](
	batchFunc BatchFunc[TKey, TValue],
	keyFunc KeyFunc[TKey, TValue],
) *dataloader.Loader[TKey, *TValue] {
	return dataloader.NewBatchedLoader(
		func(ctx context.Context, keys []TKey) []*dataloader.Result[*TValue] {
			values, err := batchFunc(ctx, keys)
			if err != nil {
				results := make([]*dataloader.Result[*TValue], len(keys))
				for i := range keys {
					results[i] = &dataloader.Result[*TValue]{Error: err}
				}

				return results
			}
			valuesMap := lo.SliceToMap(values, func(item TValue) (TKey, TValue) {
				return keyFunc(item), item
			})

			return lo.Map(keys, func(item TKey, index int) *dataloader.Result[*TValue] {
				value, ok := valuesMap[item]
				if ok {
					return &dataloader.Result[*TValue]{
						Data: &value,
					}
				}

				return &dataloader.Result[*TValue]{
					Data: nil,
				}
			})
		},
		dataloader.WithBatchCapacity[TKey, *TValue](MaxBatchSize),
		dataloader.WithWait[TKey, *TValue](WaitTime),
		dataloader.WithClearCacheOnBatch[TKey, *TValue](),
	)
}
