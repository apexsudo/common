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

func TestMany_returnsSliceForKnownKey(t *testing.T) {
	t.Parallel()

	loader := dataloaders.Many(batchItems, itemKey)
	got, err := loader.Load(context.Background(), 7)()

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, item{ID: 7, Name: "item-7"}, got[0])
}

func TestMany_returnsNilSliceForUnknownKey(t *testing.T) {
	t.Parallel()

	batch := func(_ context.Context, _ []int) ([]item, error) { return nil, nil }
	loader := dataloaders.Many(batch, itemKey)
	got, err := loader.Load(context.Background(), 99)()

	require.NoError(t, err)
	assert.Nil(t, got)
}

func TestMany_groupsMultipleValuesPerKey(t *testing.T) {
	t.Parallel()

	batch := func(_ context.Context, keys []int) ([]item, error) {
		result := make([]item, 0, len(keys)*2)
		for _, k := range keys {
			result = append(result,
				item{ID: k, Name: fmt.Sprintf("item-%d-a", k)},
				item{ID: k, Name: fmt.Sprintf("item-%d-b", k)},
			)
		}

		return result, nil
	}

	loader := dataloaders.Many(batch, itemKey)
	got, err := loader.Load(context.Background(), 5)()

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.ElementsMatch(t, []string{"item-5-a", "item-5-b"}, []string{got[0].Name, got[1].Name})
}

func TestMany_keepsValuesIsolatedAcrossKeys(t *testing.T) {
	t.Parallel()

	batch := func(_ context.Context, keys []int) ([]item, error) {
		result := make([]item, 0, len(keys)*2)
		for _, k := range keys {
			result = append(result,
				item{ID: k, Name: fmt.Sprintf("item-%d-a", k)},
				item{ID: k, Name: fmt.Sprintf("item-%d-b", k)},
			)
		}

		return result, nil
	}

	loader := dataloaders.Many(batch, itemKey)
	ctx := context.Background()

	thunk1 := loader.Load(ctx, 1)
	thunk2 := loader.Load(ctx, 2)

	got1, err1 := thunk1()
	got2, err2 := thunk2()

	require.NoError(t, err1)
	require.NoError(t, err2)
	for _, v := range got1 {
		assert.Equal(t, 1, v.ID)
	}

	for _, v := range got2 {
		assert.Equal(t, 2, v.ID)
	}
}

func TestMany_propagatesBatchError(t *testing.T) {
	t.Parallel()

	batchErr := errors.New("timeout")
	batch := func(_ context.Context, _ []int) ([]item, error) { return nil, batchErr }
	loader := dataloaders.Many(batch, itemKey)
	_, err := loader.Load(context.Background(), 1)()

	assert.ErrorIs(t, err, batchErr)
}

func TestMany_batchesMultipleLoads(t *testing.T) {
	t.Parallel()

	var callCount atomic.Int32
	batch := func(ctx context.Context, keys []int) ([]item, error) {
		callCount.Add(1)

		return batchItems(ctx, keys)
	}

	loader := dataloaders.Many(batch, itemKey)
	ctx := context.Background()

	thunk1 := loader.Load(ctx, 10)
	thunk2 := loader.Load(ctx, 20)

	items1, err1 := thunk1()
	items2, err2 := thunk2()

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.Len(t, items1, 1)
	require.Len(t, items2, 1)
	assert.Equal(t, 10, items1[0].ID)
	assert.Equal(t, 20, items2[0].ID)
	assert.Equal(t, int32(1), callCount.Load())
}
