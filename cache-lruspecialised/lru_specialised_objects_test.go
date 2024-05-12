package cachelruspecialised

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type Person struct {
	Name string
	Age  uint
}

func TestPerson_LRU(t *testing.T) {
	ttl := 10 * time.Millisecond

	cache := NewCacheLRU[int, Person](
		&ParamsNewCacheLRU{
			Capacity: 16,
			TTL:      ttl,
		},
	)

	key := 1
	person := Person{
		Name: "John",
		Age:  11,
	}

	cache.PutTTL(
		key,
		person,
	)

	cachedPerson, errGet := cache.Get(key)
	require.NoError(t, errGet)
	require.NotNil(t, cachedPerson)
	require.Equal(t,
		person,
		*cachedPerson,
	)
}

type Car struct {
	Model string
}

type Assets struct {
	Person
	Car
}

func TestAssets_LRU(t *testing.T) {
	ttl := 10 * time.Millisecond

	cache := NewCacheLRU[int, Assets](
		&ParamsNewCacheLRU{
			Capacity: 16,
			TTL:      ttl,
		},
	)

	key := 1
	person := Person{
		Name: "John",
		Age:  11,
	}

	cache.PutTTL(
		key,
		Assets{
			Person: person,
		},
	)

	cachedAsset, errGet := cache.Get(key)
	require.NoError(t, errGet)
	require.NotNil(t, cachedAsset)
	require.Equal(t,
		person,
		cachedAsset.Person,
	)
}
