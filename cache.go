package cache

import (
	"time"
)

// Result is the code of a cache Get operation
type Result string

const (
	// OK : object for key found and fresh
	OK Result = "ok"
	// NotFound : object for key not found
	NotFound Result = "not found"
	// Stale : object for key found, but stale
	Stale Result = "stale"
)

// Cache is the abstract interface for a cache that requires an TTL duration for objects
type Cache interface {

	// Put an object into the cache and expire after the given ttl
	Put(key string, value interface{}, ttl time.Duration)

	// Get an object from the cache. Stale objects are returned, but marked as stale
	Get(key string) (value interface{}, status Result)

	// Size returns the max items that may be stored in the cache
	Size() int

	// Remove removes the provided key from the cache.
	Remove(key string)
}

// New constructs a new LRU/TTL cache with the given max size (num objects)
func New(maxSize int) Cache {
	lruCache, _ := NewLRUCache(maxSize)
	return &cache{
		lruCache,
		maxSize,
	}
}

type cache struct {
	lruCache *LRUCache
	maxSize  int
}

type cacheEntry struct {
	expiry time.Time
	obj    interface{}
}

func (c *cache) Put(key string, value interface{}, ttl time.Duration) {
	c.lruCache.Add(key, cacheEntry{
		time.Now().Add(ttl),
		value,
	})
}

func (c *cache) Get(key string) (_ interface{}, _ Result) {
	entryIface, ok := c.lruCache.Get(key)
	if !ok {
		return nil, NotFound
	}

	entry := entryIface.(cacheEntry)

	if entry.expiry.After(time.Now()) {
		return entry.obj, OK
	}

	return entry.obj, Stale
}

func (c *cache) Size() int {
	return c.maxSize
}

// Remove the provided key from the cache.
func (c *cache) Remove(key string) {
	c.lruCache.Remove(key)
}
