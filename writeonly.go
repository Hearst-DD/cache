package cache

import "time"

// WriteOnly wraps a cache in a write-only wrapper (Get() always returns null)
func WriteOnly(c Cache) Cache {
	return &writeOnlyCache{
		c: c,
	}
}

type writeOnlyCache struct {
	c Cache
}

func (woc *writeOnlyCache) Put(key string, value interface{}, ttl time.Duration) {
	woc.c.Put(key, value, ttl)
}

func (woc *writeOnlyCache) Get(key string) (_ interface{}, _ Result) {
	return nil, NotFound
}

func (woc *writeOnlyCache) Size() int {
	return woc.c.Size()
}

func (woc *writeOnlyCache) Len() int {
	return woc.c.Len()
}
