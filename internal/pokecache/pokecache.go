package pokecache

import (
	"time"
	"sync"
	
)

type cacheEntry struct{
	createdAt time.Time
	val []byte
}

type Cache struct{
	mu sync.Mutex
	locs map[string]cacheEntry
}

func NewCache(t time.Duration)Cache {
	c := Cache{
		locs: make(map[string]cacheEntry),
	}
	go c.reapLoop(t)
	return c
}

func (c *Cache) Add(key string, val []byte){
	c.mu.Lock()         
	defer c.mu.Unlock()
	var newEntry cacheEntry
	newEntry.createdAt = time.Now()
	newEntry.val = val
	c.locs[key] =  newEntry
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.Lock()         
	defer c.mu.Unlock()
	if va, ok := c.locs[key]; !ok {
		return nil, false
	}else{
		return va.val, true
	}
}

func (c *Cache) reapLoop(t time.Duration){

	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for range ticker.C{

			c.mu.Lock()
			for key, val := range c.locs{
				if time.Since(val.createdAt) >= t{
					delete(c.locs, key)
				}
			}
			c.mu.Unlock()
		
	}

}