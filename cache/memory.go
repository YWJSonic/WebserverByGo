package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	memcache "github.com/patrickmn/go-cache"
)

type memoryCache struct {
	*memcache.Cache
}

func (m memoryCache) Get(key string) (interface{}, error) {
	v, ok := m.Cache.Get(key)
	if !ok {
		return nil, redis.ErrNil
	}
	return v, nil
}

func (m memoryCache) Exist(key string) (bool, error) {
	_, err := m.Get(key)
	if err != nil {
		if err != redis.ErrNil {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (m memoryCache) Set(key string, v interface{}, d time.Duration) error {
	m.Cache.Set(key, v, d)
	return nil
}

func (m memoryCache) Delete(key string) error {
	m.Cache.Delete(key)
	return nil
}

func (m memoryCache) Bool(key string) (bool, error) {
	return redis.Bool(m.Get(key))
}

func (m memoryCache) Bytes(key string) ([]byte, error) {
	return redis.Bytes(m.Get(key))
}

func (m memoryCache) Int(key string) (int, error) {
	return redis.Int(m.Get(key))
}

func (m memoryCache) Int64(key string) (int64, error) {
	return redis.Int64(m.Get(key))
}

func (m memoryCache) Float64(key string) (float64, error) {
	return redis.Float64(m.Get(key))
}

func (m memoryCache) Uint64(key string) (uint64, error) {
	return redis.Uint64(m.Get(key))
}

func (m memoryCache) String(key string) (string, error) {
	return redis.String(m.Get(key))
}

func (m memoryCache) Flush() error {
	m.Cache.Flush()
	return nil
}

func newMemoryCache(exp, dur time.Duration) Cache {
	return memoryCache{memcache.New(exp, dur)}
}
