package cachelruany

import (
	"fmt"
	"sync"
	"testing"

	"github.com/TudorHulban/caches/apperrors"
	"github.com/stretchr/testify/require"
)

func TestLRU_Delete(t *testing.T) {
	cache := NewCacheLRU(10)

	k1, v1 := 1, 77

	cache.Put(k1, v1)

	t.Log("length after first insert:", cache.Queue.Len())

	reconstructedV1, errGet1 := cache.Get(k1)
	require.NoError(t, errGet1)
	require.Equal(t, v1, reconstructedV1)

	k2, v2 := 2, 78

	cache.Put(k2, v2)

	t.Log("length after second insert:", cache.Queue.Len())

	reconstructedV2, errGet2 := cache.Get(k2)
	require.NoError(t, errGet2)
	require.Equal(t, v2, reconstructedV2)

	current := cache.Queue.Len()

	errDel := cache.Delete(k2)
	require.NoError(t, errDel)

	reconstructedV3, errGet3 := cache.Get(k2)
	require.ErrorIs(t,
		errGet3,
		apperrors.ErrRecordNotFound{},
	)
	require.Zero(t, reconstructedV3)
	require.Equal(t, current-1, cache.Queue.Len())
}

func TestLRU2(t *testing.T) {
	cache := NewCacheLRU(2)

	k1, v1 := 1, 77

	cache.Put(k1, v1)

	reconstructedV1, errGet1 := cache.Get(k1)
	require.NoError(t, errGet1)
	require.Equal(t, v1, reconstructedV1)

	k2, v2 := 2, 87

	cache.Put(k2, v2)

	reconstructedV2, errGet2 := cache.Get(k2)
	require.NoError(t, errGet2)
	require.Equal(t, v2, reconstructedV2)

	v3 := 78

	cache.Put(k1, v3)

	reconstructedV3, errGet3 := cache.Get(k1)
	require.NoError(t, errGet3)
	require.Equal(t, v3, reconstructedV3)

	fmt.Println(cache)
}

func TestLRU2Conc(_ *testing.T) {
	cache := NewCacheLRU(2)

	var wg sync.WaitGroup

	put := func(key int, value int) {
		defer wg.Done()

		cache.Put(key, value)
	}

	k1, v1 := 1, 77
	wg.Add(1)

	go put(k1, v1)

	k2, v2 := 2, 87
	wg.Add(1)

	go cache.Get(k1)
	go cache.Get(k1)
	go cache.Get(k1)

	go put(k2, v2)

	v3 := 78
	wg.Add(1)

	go put(k1, v3)

	wg.Wait()

	fmt.Println(cache)
}

func TestLRU1Any(t *testing.T) {
	cache := NewCacheLRU(1)

	k1, v1 := "1", 77

	cache.Put(k1, v1)

	reconstructedV1, errGet1 := cache.Get(k1)
	require.NoError(t, errGet1)
	require.Equal(t, v1, reconstructedV1)

	v2 := 78

	cache.Put(k1, v2)

	reconstructedV2, errGet2 := cache.Get(k1)
	require.NoError(t, errGet2)
	require.Equal(t, v2, reconstructedV2)
}

func TestLRU_Capacity1(t *testing.T) {
	cache := NewCacheLRU(1)

	k1, v1 := 1, 77

	cache.Put(k1, v1)

	reconstructedV1, errGet1 := cache.Get(k1)
	require.NoError(t, errGet1)
	require.Equal(t, v1, reconstructedV1)

	k2, v2 := 2, 78

	cache.Put(k2, v2)

	reconstructedV2, errGet2 := cache.Get(k2)
	require.NoError(t, errGet2)
	require.Equal(t, v2, reconstructedV2)
}
