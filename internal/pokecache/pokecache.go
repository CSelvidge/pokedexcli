package pokecache

import (
	"sync"
	"time"
	"errors"
)

type Cache struct {
	cache map[string]cacheEntry
	intervalTimer time.Duration //time cache is alloweed to live
	mu  *sync.RWMutex //RWMutex for concurrent access
}

type cacheEntry struct {
	createdAt time.Time
	val []byte //can be any data that can be marshaled into bytes
}

func NewCache(durationType string,  durationLife int ) (*Cache, error) {
	invalidDurationType := errors.New("invalid duration type provided")
	valueIsNil := errors.New("One or more values provided are nil")

	if durationType == "" || durationLife <= 0 {
		return nil, valueIsNil
	}

	durationMap := map[string]time.Duration{
		"second": time.Second,
		"minute":  time.Minute,
		"hour":    time.Hour, //program will not live long enough for this to work, but included for fun
	}

	timeVersion, exists := durationMap[durationType]
	if !exists {
		return nil, invalidDurationType //invalid duration type
	}

	c := &Cache{
		cache: make(map[string]cacheEntry),
		intervalTimer: time.Duration(durationLife) * timeVersion,
		mu:   &sync.RWMutex{}, 
	}

	go c.reapLoop() //start the reaping loop immediately after cache creation to begin cuncurrent reaping

	return c, nil

}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	entry, exists := c.cache[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.intervalTimer)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.RLock()
		if len(c.cache) == 0 {
			c.mu.RUnlock()
			continue
		}

		var keyDeletion []string
		for key, entry := range c.cache {
			if time.Since(entry.createdAt) > c.intervalTimer {
				keyDeletion = append(keyDeletion, key)
			}
		}
		c.mu.RUnlock()

		if len(keyDeletion) > 0 {
			c.mu.Lock()
			for _, key := range keyDeletion {
				delete(c.cache, key)
			}
			c.mu.Unlock()
		}
	}
}
