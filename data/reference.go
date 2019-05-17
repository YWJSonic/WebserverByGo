package data

import (
	"fmt"
	"time"

	"../cache"
)

// MaintainStartTime cron maintain schedule
const MaintainStartTime string = "0 0 15 * * *"

// MaintainFinishTime cron maintain schedule
const MaintainFinishTime string = "0 0 16 * * *"

// GameTypeID this server game id
const GameTypeID string = "A173D52E01A6EB65A5D6EDFB71A8C39C"

// IP Server Listen address
const IP string = "192.168.1.15"

// PORT ServerListen PORT
const PORT string = "8000"

// DBIP server connect DB address
const DBIP string = "127.0.0.1"

// DBPORT server connect DB port
const DBPORT string = "3306"

// DBUser Connect name
const DBUser string = "root"

// DBPassword connect Password
const DBPassword string = "123456"

// AccountEncodeStr account encode noise
const AccountEncodeStr string = "yrgb$"

// RedisURL cache server address
const RedisURL string = "192.168.1.15:6379"

// CacheDeleteTime cache keep time
const CacheDeleteTime time.Duration = time.Hour

// ConnectTimeOut Client connect time out
const ConnectTimeOut int64 = 15e9

// Maintain Is sow maintain time
var Maintain = false

// CacheRef old cache not used
var CacheRef *cache.Cache

// Setting settint from db data
var Setting map[string]interface{}

func init() {
	Setting = make(map[string]interface{})
	// NewCache()
}

// func NewCache() {

// 	cache, err := cache.New(cache.Config{
// 		RedisURL:  "redis://127.0.0.1:6379",
// 		MustRedis: true,
// 	})

// 	CacheRef = &cache
// 	if err != nil {
// 		panic(err)
// 	}
// }

// ServerURL ...
func ServerURL() string {
	return fmt.Sprintf("%s:%s", IP, PORT)
}
