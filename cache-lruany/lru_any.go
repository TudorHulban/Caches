package cachelruany

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"github.com/TudorHulban/caches/apperrors"
)

// Item points to a corresponding element in the doubly linked list.
type Item struct {
	KeyPtr  *list.Element
	Payload any
}

// CacheLRU acts as LRU caching based on the capacity and the cache.
// The queue is the double linked list from which the LRU element is deleted as needed.
// The cache key should be a hasheable type.
type CacheLRU struct {
	Queue *list.List
	Cache map[any]*Item // key would be the type used in correspondent method

	mu sync.Mutex

	capacity uint16
}

func NewCacheLRU(capacity uint16) *CacheLRU {
	return &CacheLRU{
		Queue:    list.New(),
		Cache:    make(map[any]*Item),
		capacity: capacity,
	}
}

// Stringer added for debugging.
func (c *CacheLRU) String() string {
	var res []string

	res = append(res, fmt.Sprintf("Capacity: %d", c.capacity)) //nolint:gocritic
	res = append(res, "Cached:")                               //nolint:gocritic

	for key, item := range c.Cache {
		res = append(res, fmt.Sprintf("key: %v, value: %v", key, item.Payload))
	}

	return strings.Join(res, "\n")
}

// Put places a key value in the cache and linked list.
// The LRU element is replaced by the added element if the capacity is reached.
func (c *CacheLRU) Put(key, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, exists := c.Cache[key]; !exists {
		if int(c.capacity) == len(c.Cache) {
			lru := c.Queue.Back()
			c.Queue.Remove(lru)

			delete(c.Cache, lru.Value.(int))
		}

		c.Cache[key] = &Item{
			KeyPtr:  c.Queue.PushFront(key),
			Payload: value,
		}
	} else {
		node.Payload = value

		c.Cache[key] = node
		c.Queue.MoveToFront(node.KeyPtr)
	}
}

// Get fetches an element from cache if it finds it by the passed key.
func (c *CacheLRU) Get(key any) (any, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, exists := c.Cache[key]; exists {
		c.Queue.MoveToFront(item.KeyPtr)

		return item.Payload, nil
	}

	return nil,
		apperrors.ErrRecordNotFound{}
}

func (c *CacheLRU) Delete(key any) error {
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
