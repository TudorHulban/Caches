package cachelruspecialised

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkSpecialisedLRU-16    	15391800	        66.34 ns/op	       0 B/op	       0 allocs/op

func BenchmarkSpecialisedLRU(b *testing.B) {
	cache := NewCacheLRU[string, int](
		&ParamsNewCacheLRU{
			Capacity: 1,
		},
	)

	key := "1"

	cache.Put(
		key,
		77,
	)

	b.ResetTimer()
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		cache.Get(key)
	}
}
