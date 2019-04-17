package cache

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	redisURL string
	rc       Cache
)

func TestMain(m *testing.M) {
	flag.StringVar(&redisURL, "r", "", "Redis server URL")
	flag.Parse()

	if redisURL == "" {
		fmt.Println("cache: Redis URL not provided")
		os.Exit(1)
	}

	c, err := New(Config{
		RedisURL:  redisURL,
		MustRedis: true,
	})
	if err != nil {
		fmt.Println("New failed:", err)
		os.Exit(1)
	}
	rc = c.(*redisCache)

	os.Exit(m.Run())
}
