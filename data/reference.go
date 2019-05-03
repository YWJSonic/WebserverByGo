package data

import (
	"fmt"
	"time"

	"../cache"
)

const IP string = "192.168.1.15"
const PORT string = "8000"

const DBIP string = "127.0.0.1"
const DBPORT string = "3306"
const DBUser string = "root"
const DBPassword string = "123456"

const DBDataSourceName string = "root:123456@/gamedb"
const DBLogSourceName string = "root:123456@/logdb"
const AccountEncodeStr string = "yrgb$"
const RedisURL string = "192.168.1.15:6379"
const CacheDeleteTime time.Duration = time.Hour

// ConnectTimeOut Client connect time out
const ConnectTimeOut int64 = 15

// Maintain Is sow maintain time
var Maintain bool

var CacheRef *cache.Cache

var Setting map[string]interface{}

func init() {
	Maintain = false
	Setting = make(map[string]interface{})
	NewCache()
}

func NewCache() {

	cache, err := cache.New(cache.Config{
		RedisURL:  "redis://127.0.0.1:6379",
		MustRedis: true,
	})

	CacheRef = &cache
	if err != nil {
		panic(err)
	}
}

// ServerURL ...
func ServerURL() string {
	return fmt.Sprintf("%s:%s", IP, PORT)
}
