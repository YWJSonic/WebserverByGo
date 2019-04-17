package cache

import (
	"math"
	"testing"
	"time"
)

func TestRedisReplace(t *testing.T) {
	defer rc.Flush()

	if err := rc.Replace("foo", "bar", NoExpire); err == nil {
		t.Error("redisCache.Replace failed: unexisting key should not be replaceable")
	}

	if err := rc.Add("foo", 2009, NoExpire); err != nil {
		t.Fatal("redisCache.Replace failed:", err)
	}

	want := 2018
	if err := rc.Replace("foo", want, NoExpire); err != nil {
		t.Fatal("redisCache.Replace failed:", err)
	}

	i, err := rc.Int("foo")
	if err != nil {
		t.Fatal("redisCache.Int failed:", err)
	}
	if i != want {
		t.Errorf("redisCache.Replace failed: Value not replaced\nGot : %d\nWant: %d", i, want)
	}
}

func TestSet(t *testing.T) {
	defer rc.Flush()

	// case 1: custom timeout, default NoExpire
	if err := rc.Set("pi", math.Pi, time.Millisecond); err != nil {
		t.Fatal("redisCache.Set failed:", err)
	}

	f, err := rc.Float64("pi")
	if err != nil {
		t.Fatal("redisCache.Float64 failed:", err)
	}
	if f != float64(math.Pi) {
		t.Errorf("redisCache.Set failed: Wrong value.\nGot : %f\nWant: %f", f, math.Pi)
	}

	time.Sleep(2 * time.Millisecond)
	if _, err = rc.Float64("pi"); err == nil {
		t.Error("redisCache.set failed: value not expired")
	}

	// case 2: use default timeout, default 1ms
	rc.(*redisCache).exp = time.Millisecond
	defer func() {
		rc.(*redisCache).exp = NoExpire
	}()

	if err = rc.Set("euler", math.E, DefaultExpire); err != nil {
		t.Fatal("redisCache.Set failed:", err)
	}

	time.Sleep(2 * time.Millisecond)
	if _, err = rc.Float64("euler"); err == nil {
		t.Error("redisCache.set failed: value not expired")
	}
}

func TestRedisDelete(t *testing.T) {
	defer rc.Flush()

	if err := rc.Add("foo", "bar", DefaultExpire); err != nil {
		t.Fatal("redisCache.Add failed:", err)
	}

	if err := rc.Delete("foo"); err != nil {
		t.Fatal("redisCache.Delete failed:", err)
	}

	_, err := rc.String("foo")
	if err == nil {
		t.Fatal("redisCache.Delete failed: key not removed")
	}
}

func TestRedisFlush(t *testing.T) {
	defer rc.Flush()

	if err := rc.Add("foo", "bar", DefaultExpire); err != nil {
		t.Fatal("redisCache.Add failed:", err)
	}

	if err := rc.Flush(); err != nil {
		t.Fatal("redisCache.Flush failed:", err)
	}

	_, err := rc.String("foo")
	if err == nil {
		t.Fatal("redisCache.Flush failed: All keys should be deleted")
	}
}
