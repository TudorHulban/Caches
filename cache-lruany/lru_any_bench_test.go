package cachelruany

import "testing"

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkSumLoop-16    	 7855671	       141.6 ns/op	       0 B/op	       0 allocs/op

func BenchmarkAnyLRU(b *testing.B) {
	cache := NewCacheLRU(1)

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
