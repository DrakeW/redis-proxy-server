package cache

import (
	"fmt"
	"sync"
	"time"
)

// Config represents the config options for LRU cache
type Config struct {
	Expiry     time.Duration
	MaxEntries uint
}

// LRU represents a LRU (least recently used) Cache
type LRU struct {
	content map[string]*entry
	lock    *sync.Mutex
	config  *Config
}

type entry struct {
	val          interface{}
	latestAccess time.Time
	lock         *sync.Mutex
}

// NewLRUCache - returns a LRU cache object
func NewLRUCache(config *Config) *LRU {
	return &LRU{
		content: make(map[string]*entry),
		lock:    new(sync.Mutex),
		config:  config,
	}
}

// Get retrieves value from cache based on key and returns error if key doesn't exist or has expired
func (c *LRU) Get(key string) (interface{}, error) {
	en, ok := c.content[key]
	if !ok {
		return nil, nil
	}
	// check if the entry being retrieved has already expired
	if time.Since(en.latestAccess) >= c.config.Expiry {
		err := c.Remove(key)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("entry with key \"%s\" has expired", key)
	}
	return en.getVal(), nil
}

// Add add key-value pair into the cache based on input, if max capacity is reached, the
// least recently used entry will be expired to empty out slot for the new one.
// This operation is atomic at the cache level
func (c *LRU) Add(key string, val string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	en, ok := c.content[key]
	if !ok {
		if uint(len(c.content)) >= c.config.MaxEntries {
			// TODO: evict lRU entry
		}
		c.content[key] = &entry{
			val:          val,
			latestAccess: time.Now(),
			lock:         new(sync.Mutex),
		}
	} else {
		en.update(val)
	}
	return nil
}

// Remove delets a key-value pair from the cache
// this operation is atomic at the cache level
func (c *LRU) Remove(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.content[key]
	if !ok {
		return fmt.Errorf("failed to delete key \"%s\"", key)
	}
	delete(c.content, key)
	return nil
}

// update - updates the value of a cache entry, this operation is atomic for each entry
func (en *entry) update(val string) {
	en.lock.Lock()
	defer en.lock.Unlock()

	en.val = val
	en.latestAccess = time.Now()
}

// getVal - get the value of a cache entry, and refreshes its latest access time for each entry
// this operation is atomic
func (en *entry) getVal() interface{} {
	en.lock.Lock()
	defer en.lock.Unlock()

	en.latestAccess = time.Now()
	return en.val
}
