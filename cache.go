package bunny

import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	m map[string][]byte
}

func Open() *Cache {
	return &Cache{
		m: make(map[string][]byte),
	}
}

func (c *Cache) Get(key string) []byte {
	c.RLock()
	defer c.RUnlock()
	if data, ok := c.m[key]; ok {
		return data
	}
	return nil
}

func (c *Cache) Set(key string, value []byte) {
	c.Lock()
	defer c.Unlock()
	c.m[key] = value

}

func (c *Cache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.m, key)
}
