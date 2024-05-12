package cachettl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/require"
)

const _domain = DomainModel("xxx")

func TestIsExpired(t *testing.T) {
	cache := NewCache(context.Background(), _domain)
	defer cache.Close()

	now := time.Now()
	require.False(t, cache.isTimeExpired(now.UnixNano(), 1, now))
	require.True(t,
		cache.isTimeExpired(now.UnixNano(), 1, now.Add(time.Second*time.Duration(2))),
	)
}

func TestInMemoryCache(t *testing.T) {
	require := require.New(t)

	cache := NewCache(context.Background(), _domain)
	defer cache.Close()

	dto := NewDTO(_domain, []byte("1"))
	require.NoError(cache.Set(dto))

	fetchedDTO, errGet1 := cache.Get(dto.Key)
	t.Log(*fetchedDTO)
	t.Log(dto.Key)

	require.NoError(errGet1)
	require.Nil(
		deep.Equal(dto, fetchedDTO),
		fmt.Sprintf("dto: %#v\n", *dto),
	)

	errDel := cache.Delete(dto.Key)
	require.NoError(errDel)

	shouldBeNilDTO, errGet2 := cache.Get(dto.Key)
	require.Error(errGet2)
	require.Nil(shouldBeNilDTO)
}

func TestInMemoryCacheTTL(t *testing.T) {
	require := require.New(t)

	cache := NewCache(context.Background(), _domain, WithTTL(1), WithSecondsBetweenCleanUps(1))
	defer cache.Close()

	dto := NewDTO(_domain, []byte("1"))
	require.NoError(cache.Set(dto))

	fetchedDTO, errGet1 := cache.Get(dto.Key)
	t.Log(*fetchedDTO)
	t.Log(dto.Key)

	require.NoError(errGet1)
	require.Nil(deep.Equal(dto, fetchedDTO), fmt.Sprintf("dto: %#v\n", *dto))

	time.Sleep(1100 * time.Millisecond)

	shouldBeNilDTO, errGet2 := cache.Get(dto.Key)
	require.Error(errGet2)
	require.Nil(shouldBeNilDTO)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkCache-16    	 3039609	       367.0 ns/op	     207 B/op	       3 allocs/op
func BenchmarkCache(b *testing.B) {
	cache := NewCache(context.Background(), _domain, WithTTL(0), WithSecondsBetweenCleanUps(1))
	defer cache.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dto := NewDTO(_domain, []byte{byte(i)})

		cache.Set(dto)
		cache.Get(dto.Key)
	}
}
