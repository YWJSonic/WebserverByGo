package mycache

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"github.com/gomodule/redigo/redis"
)

var CachePool *redis.Pool

func init() {
	newCachePool()
}

const ConnectTimeout time.Duration = 20 * time.Second
const ReadTimeout time.Duration = 5 * time.Second
const WriteTimeout time.Duration = 10 * time.Second

func newCachePool() {
	CachePool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		MaxActive:   20,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", data.RedisURL,
				redis.DialConnectTimeout(ConnectTimeout),
				redis.DialReadTimeout(ReadTimeout),
				redis.DialWriteTimeout(WriteTimeout))
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			//验证redis密码
			// if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
			// 	return nil, fmt.Errorf("redis auth password error: %s", authErr)
			// }
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				return fmt.Errorf("ping redis error: %s", err)
			}
			return nil
		},
	}
}

func get(ket string) (interface{}, error) {
	Conn := CachePool.Get()
	defer Conn.Close()

	return Conn.Do("GET", ket)
}

func set(args []interface{}, d time.Duration) (interface{}, error) {
	Conn := CachePool.Get()
	defer Conn.Close()

	if d != 0 {
		args = append(args, "PX", int(d/time.Millisecond))
	}
	return Conn.Do("SET", args...)
}

func del(key string) error {
	Conn := CachePool.Get()
	defer Conn.Close()

	_, err := Conn.Do("DEL", key)
	return err
}

func runExist(key string) (bool, error) {
	Conn := CachePool.Get()
	defer Conn.Close()

	if _, err := Conn.Do("EXISTS", key); err != redis.ErrNil {
		return false, err
	} else {
		return false, nil
	}
	return true, nil
}

func runAdd(key string, v interface{}, d time.Duration) error {

	args := []interface{}{key, v, "NX"}

	_, err := set(args, d)
	return err
}

func runReplace(key string, v interface{}, d time.Duration) error {
	args := []interface{}{key, v, "XX"}

	v, err := set(args, d)
	if err != nil {
		return err
	}
	if v == nil {
		return redis.ErrNil
	}
	return nil
}

func runSet(key string, v interface{}, d time.Duration) error {
	args := []interface{}{key, v}

	_, err := set(args, d)
	if err != nil {
		errorlog.ErrorLogPrintln("Cache Set", key, err)
	}
	return err
}

func runDelete(key string) error {
	Conn := CachePool.Get()
	defer Conn.Close()

	_, err := Conn.Do("DEL", key)
	return err
}

func getBool(key string) (bool, error) {
	return redis.Bool(get(key))
}

func getBytes(key string) ([]byte, error) {
	return redis.Bytes(get(key))
}

func getInt(key string) (int, error) {
	return redis.Int(get(key))
}

func getInt64(key string) (int64, error) {
	return redis.Int64(get(key))
}

func getFloat64(key string) (float64, error) {
	return redis.Float64(get(key))
}

func getUint64(key string) (uint64, error) {
	return redis.Uint64(get(key))
}

func getString(key string) (string, error) {
	return redis.String(get(key))
}

func runFlush() error {
	Conn := CachePool.Get()
	defer Conn.Close()
	_, err := Conn.Do("FLUSHDB")
	return err
}
