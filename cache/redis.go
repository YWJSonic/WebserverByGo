package cache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type redisCache struct {
	redis.Conn
	exp time.Duration
}

func (r *redisCache) Get(key string) (interface{}, error) {
	return r.Do("GET", key)
}

func (r *redisCache) set(args []interface{}, d time.Duration) (interface{}, error) {
	if d != NoExpire {
		if d == DefaultExpire {
			if r.exp == NoExpire {
				return r.Do("SET", args...)
			}
			d = r.exp
		}
		args = append(args, "PX", int(d/time.Millisecond))
	}
	return r.Do("SET", args...)
}

func (r *redisCache) Exist(key string) (bool, error) {
	if _, err := r.Do("EXISTS", key); err != redis.ErrNil {
		return false, err
	} else {
		return false, nil
	}
	return true, nil
}

func (r *redisCache) Add(key string, v interface{}, d time.Duration) error {
	args := []interface{}{key, v, "NX"}

	_, err := r.set(args, d)
	return err
}

func (r *redisCache) Replace(key string, v interface{}, d time.Duration) error {
	args := []interface{}{key, v, "XX"}

	v, err := r.set(args, d)
	if err != nil {
		return err
	}
	if v == nil {
		return redis.ErrNil
	}
	return nil
}

func (r *redisCache) Set(key string, v interface{}, d time.Duration) error {
	args := []interface{}{key, v}

	_, err := r.set(args, d)
	return err
}

func (r *redisCache) Delete(key string) error {
	_, err := r.Do("DEL", key)
	return err
}

func (r *redisCache) Bool(key string) (bool, error) {
	return redis.Bool(r.Get(key))
}

func (r *redisCache) Bytes(key string) ([]byte, error) {
	return redis.Bytes(r.Get(key))
}

func (r *redisCache) Int(key string) (int, error) {
	return redis.Int(r.Get(key))
}

func (r *redisCache) Int64(key string) (int64, error) {
	return redis.Int64(r.Get(key))
}

func (r *redisCache) Float64(key string) (float64, error) {
	return redis.Float64(r.Get(key))
}

func (r *redisCache) Uint64(key string) (uint64, error) {
	return redis.Uint64(r.Get(key))
}

func (r *redisCache) String(key string) (string, error) {
	return redis.String(r.Get(key))
}

func (r *redisCache) Flush() error {
	_, err := r.Do("FLUSHDB")
	return err
}

func newRedisCache(c redis.Conn, exp time.Duration) Cache {
	if exp == DefaultExpire {
		exp = NoExpire
	}
	return &redisCache{c, exp}
}
