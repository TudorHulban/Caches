package cachelruspecialised

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/TudorHulban/caches/apperrors"
)

type item[V any] struct {
	keyPtr *list.Element // holds key of item in cache.

	payload      V
	timestampTTL int64
}

type CacheLRU[K comparable, V any] struct {
	Queue *list.List
	Cache map[K]*item[V]

	mu sync.Mutex

	ttl      time.Duration
	capacity uint16
}

type ParamsNewCacheLRU struct {
	TTL      time.Duration
	Capacity uint16
}

func NewCacheLRU[K comparable, V any](params *ParamsNewCacheLRU) *CacheLRU[K, V] {
	return &CacheLRU[K, V]{
		Queue: list.New(),
		Cache: make(map[K]*item[V]),

		capacity: params.Capacity,
		ttl:      params.TTL,
	}
}
func (c *CacheLRU[K, V]) String() string {
	var res []string

	res = append(res, fmt.Sprintf("Capacity: %d", c.capacity)) //nolint:gocritic
	res = append(res, "Cached:")                               //nolint:gocritic

	for key, item := range c.Cache {
		res = append(res,
			fmt.Sprintf(
				"key: %v, value: %v",
				key,
				item.payload,
			),
		)
	}

	return strings.Join(res, "\n")
}

func (c *CacheLRU[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.Cache[key]; !exists {
		if int(c.capacity) == len(c.Cache) {
			elementLRU := c.Queue.Back()
			c.Queue.Remove(elementLRU)

			delete(
				c.Cache,
				elementLRU.Value.(K),
			)
		}

		c.Cache[key] = &item[V]{
			keyPtr:  c.Queue.PushFront(key),
			payload: value,
		}

		return
	} else {
		node.payload = value

		c.Cache[key] = node
		c.Queue.MoveToFront(node.keyPtr)
	}
}

func (c *CacheLRU[K, V]) PutTTL(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.Cache[key]; !exists {
		if int(c.capacity) == len(c.Cache) {
			elementLRU := c.Queue.Back()
			c.Queue.Remove(elementLRU)

			delete(
				c.Cache,
				elementLRU.Value.(K),
			)
		}

		c.Cache[key] = &item[V]{
			keyPtr:  c.Queue.PushFront(key),
			payload: value,

			timestampTTL: time.Now().
				Add(c.ttl).
				UnixNano(),
		}

		return
	} else {
		node.payload = value
		node.timestampTTL = time.Now().
			Add(c.ttl).
			UnixNano()

		c.Cache[key] = node
		c.Queue.MoveToFront(node.keyPtr)
	}
}

func (c *CacheLRU[K, V]) Get(key K) (*V, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, exists := c.Cache[key]; exists {
		if item.timestampTTL > 0 && time.Now().UnixNano() >= item.timestampTTL {
			go c.Delete(key)

			return nil,
				apperrors.ErrRecordNotFound{}

		}

		c.Queue.MoveToFront(item.keyPtr)

		return &item.payload,
			nil
	}

	return nil,
		apperrors.ErrRecordNotFound{}
}

func (c *CacheLRU[K, V]) Delete(key K) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.Cache[key]; exists {
		currentNode := c.Queue.Back()

		removeIfFound := func(node *list.Element) bool {
			if node.Value == key {
				c.Queue.Remove(node)

				delete(c.Cache, key)

				return true
			}

			return false
		}

		for {
			if removeIfFound(currentNode) {
				return nil
			}

			currentNode = currentNode.Prev()

			if currentNode.Prev() == nil {
				if removeIfFound(currentNode) {
					return nil
				}

				break
			}
		}
	}

	return apperrors.ErrRecordNotFound{}
}
