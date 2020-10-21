package appcache

import (
	"encoding/json"
	"time"

	"sync"

	"github.com/patrickmn/go-cache"
)

type InMemoryCache struct {
	cache cache.Cache
	mu sync.RWMutex
	defaultExpiration time.Duration
}

func (c InMemoryCache) Get(key string, ptrValue interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, found := c.cache.Get(key)
	if !found {
		return ErrCacheMiss
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, ptrValue)
}

func (c InMemoryCache) Set(key string, value interface{}, expires time.Duration) error {
	panic("implement me")
}

func (c InMemoryCache) SetFields(key string, value map[string]interface{}, expires time.Duration) {
	panic("implement me")
}

func (c InMemoryCache) GetMulti(keys ...string) (Getter, error) {
	panic("implement me")
}

func (c InMemoryCache) Delete(key string) error {
	panic("implement me")
}

func (c InMemoryCache) Replace(key string, value interface{}, expires time.Duration) error {
	panic("implement me")
}

func (c InMemoryCache) Flush() error {
	panic("implement me")
}

func (c InMemoryCache) Keys() ([]string, error) {
	panic("implement me")
}

func NewInMemoryCache(defaultExpiration time.Duration) InMemoryCache {
	return InMemoryCache{
		cache:             *cache.New(defaultExpiration, time.Minute),
		mu:                sync.RWMutex{},
		defaultExpiration: defaultExpiration,
	}
}

func (c InMemoryCache) Add(key string, value interface{}, expires time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.cache.Add(key, value, expires)
	if err != nil {
		return ErrNotStored
	}
	return err
}
