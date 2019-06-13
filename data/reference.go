package data

import (
	"fmt"
	"time"
)

// CacheDeleteTime cache keep time
const CacheDeleteTime time.Duration = time.Hour

// ConnectTimeOut Client connect time out
const ConnectTimeOut int64 = 15e9

// MaintainStartTime cron maintain schedule
var MaintainStartTime = "0 0 15 * * *"

// MaintainFinishTime cron maintain schedule
var MaintainFinishTime = "0 0 16 * * *"

// GameTypeID this server game id
var GameTypeID = "A173D52E01A6EB65A5D6EDFB71A8C39C"

// IP Server Listen address
var IP = "127.0.0.1"

// PORT ServerListen PORT
var PORT = "8000"

// DBIP server connect DB address
var DBIP = "127.0.0.1"

// DBPORT server connect DB port
var DBPORT = "3306"

// DBUser Connect name
var DBUser = "serverConnect"

// DBPassword connect Password
var DBPassword = "123qweasdzxc"

// AccountEncodeStr account encode noise
var AccountEncodeStr = "yrgb$"

// RedisURL cache server address
var RedisURL = "127.0.0.1:6379"

// Maintain Is sow maintain time
var Maintain = false

// GameResultURL gamelogic server API URL
// const GameResultURL = "http://192.168.1.146:9781/api/entry"

// CacheRef old cache not used
// var CacheRef *cache.Cache

// Setting settint from db data
var Setting map[string]interface{}

func init() {
	Setting = make(map[string]interface{})
	// NewCache()
}

// ServerURL ...
func ServerURL() string {
	return fmt.Sprintf("%s:%s", IP, PORT)
}
