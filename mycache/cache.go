package mycache

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

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
	}
	return false, nil

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
