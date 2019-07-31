package cache

import (
	"testing"
	"time"
)

func TestLRUCache(t *testing.T) {
	cache := NewLRUCache(&Config{
		Expiry:     100 * time.Millisecond,
		MaxEntries: 2,
	})

	t.Run("Add - normal", func(t *testing.T) {
		cache.Add("a", 1)
		res := cache.Get("a").(int)
		if res != 1 {
			t.Errorf("wrong value stored in cache. got: %d, expected: %d", res, 1)
		}
	})

	t.Run("Add - with eviction", func(t *testing.T) {
		cache.Add("b", 2)
		cache.Add("c", 3)
		res := cache.Get("c").(int)
		if res != 3 {
			t.Errorf("wrong value stored in cache. got: %d, expected: %d", res, 3)
		}
		res2 := cache.Get("a")
		if res2 != nil {
			t.Errorf("LRU cache entry with key=%s didn't get evicted", "a")
		}
	})

	t.Run("Get - cache hit", func(t *testing.T) {
		res := cache.Get("b").(int)
		if res != 2 {
			t.Errorf("wrong value stored in cache. got: %d, expected: %d", res, 2)
		}
	})

	t.Run("Get - cache expire", func(t *testing.T) {
		cache.Add("should-expire", 4)
		time.Sleep(150 * time.Millisecond)
		res := cache.Get("should-expire")
		if res != nil {
			t.Errorf("cache entry should have expired. got: %d, expected: %v", res, nil)
		}
	})

	t.Run("Get - cache miss", func(t *testing.T) {
		res := cache.Get("not-there")
		if res != nil {
			t.Errorf("should be a cach miss. got: %d, expected: %v", res, nil)
		}
	})

	t.Run("Remove", func(t *testing.T) {
		cache.Remove("a")
		res := cache.Get("a")
		if res != nil {
			t.Errorf("cache entry should be removed for key %s. got: %d, expected: %v", "a", res, nil)
		}
	})

}
