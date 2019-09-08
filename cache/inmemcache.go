package cache

import (
	"errors"
	"sync"
)

type inMemCache struct {
	c    map[string][]byte
	lock sync.RWMutex
	Stat
}

func (c *inMemCache) Set(k string, v []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if tmp, ok := c.c[k]; ok {
		c.del(k, tmp)
	}

	c.c[k] = v
	c.add(k, v)

	return nil
}

func (c *inMemCache) Get(k string) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if tmp, ok := c.c[k]; ok {
		return tmp, nil
	}
	return nil, errors.New(ErrKeyNotFound)
}

func (c *inMemCache) Del(k string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if tmp, ok := c.c[k]; ok {
		c.del(k, tmp)
		delete(c.c, k)
	}
	return nil
}

func (c *inMemCache) GetStat() Stat {
	return c.Stat
}

func newInMemCache() *inMemCache {
	return &inMemCache{
		c:    make(map[string][]byte),
		lock: sync.RWMutex{},
		Stat: Stat{},
	}
}
