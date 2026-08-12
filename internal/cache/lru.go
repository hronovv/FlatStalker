package cache

import (
	"container/list"
	"sync"
	"time"
)

const (
	DefaultSize = 1024
	DefaultTTL  = 5 * time.Minute
)

// LRU is a small thread-safe LRU cache with per-entry TTL.
type LRU[V any] struct {
	mu       sync.Mutex
	capacity int
	ttl      time.Duration
	items    map[int64]*list.Element
	order    *list.List
}

type entry[V any] struct {
	key       int64
	value     V
	expiresAt time.Time
}

func NewLRU[V any](capacity int, ttl time.Duration) *LRU[V] {
	if capacity < 1 {
		capacity = 1
	}
	return &LRU[V]{
		capacity: capacity,
		ttl:      ttl,
		items:    make(map[int64]*list.Element, capacity),
		order:    list.New(),
	}
}

func (c *LRU[V]) Get(key int64) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	el, ok := c.items[key]
	if !ok {
		return zero, false
	}
	e := el.Value.(*entry[V])
	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		c.removeElement(el)
		return zero, false
	}
	c.order.MoveToFront(el)
	return e.value, true
}

func (c *LRU[V]) Put(key int64, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiresAt := time.Time{}
	if c.ttl > 0 {
		expiresAt = time.Now().Add(c.ttl)
	}

	if el, ok := c.items[key]; ok {
		e := el.Value.(*entry[V])
		e.value = value
		e.expiresAt = expiresAt
		c.order.MoveToFront(el)
		return
	}

	for c.order.Len() >= c.capacity {
		c.removeElement(c.order.Back())
	}

	el := c.order.PushFront(&entry[V]{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
	})
	c.items[key] = el
}

func (c *LRU[V]) Delete(key int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.removeElement(el)
	}
}

func (c *LRU[V]) removeElement(el *list.Element) {
	if el == nil {
		return
	}
	e := el.Value.(*entry[V])
	delete(c.items, e.key)
	c.order.Remove(el)
}
