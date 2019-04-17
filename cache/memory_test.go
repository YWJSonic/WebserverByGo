package cache

import "testing"

func TestMemGet(t *testing.T) {
	mc, err := New(Config{
		DefaultExpiration: NoExpire,
		CleanupInterval:   NoCleanupInterval,
	})
	if err != nil {
		t.Fatal("New failed:", err)
	}

	if err = mc.Add("foo", "bar", DefaultExpire); err != nil {
		t.Fatal("memoryCache.Add failed:", err)
	}

	if _, err = mc.Get("baz"); err == nil {
		t.Error("memoryCache.Get failed: getting non-set key should return error")
	}

	v, err := mc.Get("foo")
	if err != nil {
		t.Error("memoryCache.Get failed:", err)
	}
	if s := v.(string); s != "bar" {
		t.Errorf("memoryCache.Get failed: Wrong value returned\nGot : %s\nWant: %v", s, v)
	}
}
