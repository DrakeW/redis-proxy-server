package cache

import (
	"container/list"
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
	entryMap map[string]*list.Element
	lock     *sync.Mutex
	config   *Config
	content  *list.List
}

type entry struct {
	key          string
	val          interface{}
	latestAccess time.Time
}

// NewLRUCache - returns a LRU cache object
func NewLRUCache(config *Config) *LRU {
	return &LRU{
		entryMap: make(map[string]*list.Element),
		lock:     new(sync.Mutex),
		config:   config,
		content:  list.New(),
	}
}

// Get retrieves value from cache based on key and returns error if key doesn't exist or has expired
func (c *LRU) Get(key string) interface{} {
	elem, ok := c.entryMap[key]
	if !ok {
		return nil
	}
	en := elem.Value.(*entry)
	// check if the entry being retrieved has already expired
	if time.Since(en.latestAccess) >= c.config.Expiry {
		c.Remove(key)
		return nil
	}
	// get value and update cache metadata
	defer c.refreshEntryMetadata(elem)
	return en.val
}

// Add add key-value pair into the cache based on input, if max capacity is reached, the
// least recently used entry will be expired to empty out slot for the new one.
// This operation is atomic at the cache level
func (c *LRU) Add(key string, val string) error {
	elem, ok := c.entryMap[key]
	if !ok {
		if uint(c.content.Len()) >= c.config.MaxEntries {
			lruEntry := c.content.Back().Value.(*entry)
			// remove the entry at the tail of the list (LRU)
			c.Remove(lruEntry.key)
		}

		c.lock.Lock()
		defer c.lock.Unlock()

		newEntry := &entry{
			key:          key,
			val:          val,
			latestAccess: time.Now(),
		}
		newElem := c.content.PushFront(newEntry)
		c.entryMap[key] = newElem
	} else {
		c.update(elem, val)
	}
	return nil
}

// Remove delets a key-value pair from the cache
// this operation is atomic at the cache level
func (c *LRU) Remove(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	elem, ok := c.entryMap[key]
	if ok {
		delete(c.entryMap, key)
		c.content.Remove(elem)
	}
}

// update an existing cache entry's value
func (c *LRU) update(elem *list.Element, val interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()

	en := elem.Value.(*entry)
	en.val = val
	defer c.refreshEntryMetadata(elem)
}

// Refresh the entry's latest access time and update its position in the cache content
// This function should be called after an existing entry is accessed
func (c *LRU) refreshEntryMetadata(elem *list.Element) {
	c.lock.Lock()
	defer c.lock.Unlock()

	en := elem.Value.(*entry)
	en.latestAccess = time.Now()
	c.content.MoveToFront(elem)
}
