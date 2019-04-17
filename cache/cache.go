// Package cache wraps the package "github.com/patrickmn/go-cache"
// and "github.com/gomodule/redigo/redis". It provides universal interface
// which makes it easy to switch between in-memory cache and Redis server.
package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
	memcache "github.com/patrickmn/go-cache"
)

const (
	// NoCleanupInterval indicates that never drop expired keys
	// if using in-memory cache.
	NoCleanupInterval = 0

	// NoExpire indicates that certain key never expires.
	NoExpire = memcache.NoExpiration

	// DefaultExpire indicates that the cache should apply default
	// expiration time to certain key.
	DefaultExpire = memcache.DefaultExpiration
)

// A Cache defines the common methods of a cache.
type Cache interface {
	// Get returns the value of given key, returning error if the
	// key doesn't exist.
	Get(key string) (interface{}, error)

	Exist(key string) (bool, error)

	// Add applies a value if the key doesn't exist.
	// Else it returns error.
	Add(key string, v interface{}, d time.Duration) error

	// Replace updates the value only if the key exists.
	// Else it returns error.
	Replace(key string, v interface{}, d time.Duration) error

	// Set updates the value, adding it if the key doesn't exist.
	Set(key string, v interface{}, d time.Duration) error

	// Delete removes the value, doing nothing if the key doesn't exist.
	Delete(key string) error

	Bool(key string) (bool, error)
	Bytes(key string) ([]byte, error)
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Float64(key string) (float64, error)
	Uint64(key string) (uint64, error)
	String(key string) (string, error)

	// Flush deletes all the keys from the cache.
	Flush() error
}

// A Config spcefies the options to create a cache.
type Config struct {
	DefaultExpiration time.Duration
	CleanupInterval   time.Duration
	RedisURL          string
	RedisOptions      []redis.DialOption
	MustRedis         bool
}

// New returns a new initialized cache with given config.
func New(cfg Config) (Cache, error) {
	if cfg.MustRedis {
		c, err := redis.DialURL(cfg.RedisURL, cfg.RedisOptions...)
		if err != nil {
			return nil, err
		}

		return newRedisCache(c, cfg.DefaultExpiration), nil
	}
	return newMemoryCache(cfg.DefaultExpiration, cfg.CleanupInterval), nil
}
