package dataloaders_test

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apexsudo/common/pkg/graphql/dataloaders"
)

type item struct {
	ID   int
	Name string
}

func itemKey(v item) int { return v.ID }

func batchItems(_ context.Context, keys []int) ([]item, error) {
	result := make([]item, 0, len(keys))
	for _, k := range keys {
		result = append(result, item{ID: k, Name: fmt.Sprintf("item-%d", k)})
	}

	return result, nil
}

func TestSingle_returnsPointerForKnownKey(t *testing.T) {
	t.Parallel()

	loader := dataloaders.Single(batchItems, itemKey)
	got, err := loader.Load(context.Background(), 42)()

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, item{ID: 42, Name: "item-42"}, *got)
}

func TestSingle_returnsNilPointerForUnknownKey(t *testing.T) {
	t.Parallel()

	batch := func(_ context.Context, _ []int) ([]item, error) { return nil, nil }
	loader := dataloaders.Single(batch, itemKey)
	got, err := loader.Load(context.Background(), 99)()

	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestSingle_propagatesBatchError(t *testing.T) {
	t.Parallel()

	batchErr := errors.New("database unavailable")
	batch := func(_ context.Context, _ []int) ([]item, error) { return nil, batchErr }
	loader := dataloaders.Single(batch, itemKey)
	_, err := loader.Load(context.Background(), 1)()

	assert.ErrorIs(t, err, batchErr)
}

func TestSingle_batchesMultipleLoads(t *testing.T) {
	t.Parallel()

	var callCount atomic.Int32
	batch := func(ctx context.Context, keys []int) ([]item, error) {
		callCount.Add(1)

		return batchItems(ctx, keys)
	}

	loader := dataloaders.Single(batch, itemKey)
	ctx := context.Background()

	// Register all loads before resolving any, so they land in one batch.
	thunk1 := loader.Load(ctx, 1)
	thunk2 := loader.Load(ctx, 2)
	thunk3 := loader.Load(ctx, 3)

	item1, err1 := thunk1()
	item2, err2 := thunk2()
	item3, err3 := thunk3()

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)
	assert.Equal(t, 1, item1.ID)
	assert.Equal(t, 2, item2.ID)
	assert.Equal(t, 3, item3.ID)
	assert.Equal(t, int32(1), callCount.Load())
}
