package cache

import (
	"container/list"
	"errors"
	"sync"
)

// Cache is a thread-safe fixed size LRU cache.
type LRUCache struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	lock      sync.Mutex
}

// entry is used to hold a value in the evictList
type entry struct {
	key   interface{}
	value interface{}
}

// NewLRUCache creates an LRU of the given size
func NewLRUCache(size int) (*LRUCache, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRUCache{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element, size),
	}
	return c, nil
}

// Purge is used to completely clear the cache
func (c *LRUCache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.evictList = list.New()
	c.items = make(map[interface{}]*list.Element, c.size)
}

// Add adds a value to the cache.
func (c *LRUCache) Add(key, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		return
	}

	// Add new item
	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	// Verify size not exceeded
	if c.evictList.Len() > c.size {
		c.removeOldest()
	}
}

// Get looks up a key's value from the cache.
func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache.
func (c *LRUCache) Remove(key interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRUCache) RemoveOldest() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.removeOldest()
}

// removeOldest removes the oldest item from the cache.
func (c *LRUCache) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRUCache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
}

// Len returns the number of items in the cache.
func (c *LRUCache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.evictList.Len()
}
