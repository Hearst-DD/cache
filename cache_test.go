package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldEvictLRUItemsAboveAMaxSize(t *testing.T) {
	cache := New(2)
	cache.Put("foo", "FOO", time.Minute)
	cache.Put("bar", "BAR", time.Minute)

	_, status := cache.Get("foo")
	if status != OK {
		t.Errorf("expected OK but got %s", status)
	}

	// ensure that foo is the LRU
	cache.Get("bar")

	// third add should exceed max size and cause an eviction
	cache.Put("baz", "BAZ", time.Minute)

	_, status = cache.Get("foo")
	if status != NotFound {
		t.Errorf("expected NotFound but got %s", status)
	}
	if cache.Size() != 2 {
		t.Errorf("expected cache size 2 but got %d", cache.Size())
	}
}

func TestShouldStoreObjectsInCacheAndRetrieveWithCacheKey(t *testing.T) {
	cache := New(100)

	cache.Put("foo", "FOO", time.Minute)
	obj, status := cache.Get("foo")
	if status != OK {
		t.Errorf("expected OK but got %s", status)
	}
	if obj.(string) != "FOO" {
		t.Errorf("expected FOO but got %v", obj)
	}
}

func TestShouldApplyTTLWhenInsertingObjects(t *testing.T) {
	cache := New(100)
	ttl := time.Millisecond * 10

	cache.Put("foo", "FOO", ttl)
	_, status := cache.Get("foo")
	if status != OK {
		t.Errorf("expected OK but got %s", status)
	}

	time.Sleep(time.Millisecond * 10)

	_, status = cache.Get("foo")
	if status != Stale {
		t.Errorf("expected Stale but got %s", status)
	}
}

func TestShouldReturnStaleItemAndIndicateStale(t *testing.T) {
	cache := New(100)
	ttl := time.Millisecond * 10

	cache.Put("foo", "FOO", ttl)

	time.Sleep(time.Millisecond * 10)

	obj, status := cache.Get("foo")
	if status != Stale {
		t.Errorf("expected Stale but got %s", status)
	}
	if obj.(string) != "FOO" {
		t.Errorf("expected stale object FOO back but got %v", obj)
	}
}

func TestNotFoundShouldReturnBlankItemAndIndicateNotFound(t *testing.T) {
	cache := New(100)

	_, status := cache.Get("foo")
	if status != NotFound {
		t.Errorf("expected NotFound but got %s", status)
	}
}

func Test_WriteOnly(t *testing.T) {
	cache := New(100)
	writeOnlyCache := WriteOnly(cache)

	result, ok := writeOnlyCache.Get("foo")
	assert.Nil(t, result)
	assert.Equal(t, NotFound, ok)

	writeOnlyCache.Put("foo", "bar", time.Hour)

	result, ok = writeOnlyCache.Get("foo")
	assert.Nil(t, result)
	assert.Equal(t, NotFound, ok)

	result, ok = cache.Get("foo")
	assert.Equal(t, "bar", result)
	assert.Equal(t, OK, ok)
}
