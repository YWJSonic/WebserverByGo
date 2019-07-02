package serversetting

import (
	"sync"
	"time"
)

var mu *sync.RWMutex

// CacheDeleteTime cache keep time
const CacheDeleteTime time.Duration = time.Hour

// ConnectTimeOut Client connect time out
const ConnectTimeOut int64 = 15e9

// platform api url
var (
	// MaintainStartTime cron maintain schedule
	MaintainStartTime = "0 0 15 * * *"
	// MaintainFinishTime cron maintain schedule
	MaintainFinishTime = "0 0 16 * * *"
	// GameTypeID this server game id
	GameTypeID = "A173D52E01A6EB65A5D6EDFB71A8C39C"
	// IP Server Listen address
	IP = "127.0.0.1"
	// PORT ServerListen PORT
	PORT = "8000"
	// DBIP server connect DB address
	DBIP = "127.0.0.1"
	// DBPORT server connect DB port
	DBPORT = "3306"
	// DBUser Connect name
	DBUser = "serverConnect"
	// DBPassword connect Password
	DBPassword = "123qweasdzxc"
	// AccountEncodeStr account encode noise
	AccountEncodeStr = "yrgb$"
	// RedisURL cache server address
	RedisURL = "127.0.0.1:6379"
)

// Maintain Is sow maintain time
var maintain = false

// GameResultURL gamelogic server API URL
// const GameResultURL = "http://192.168.1.146:9781/api/entry"

// CacheRef old cache not used
// var CacheRef *cache.Cache

// Setting settint from db data
var Setting map[string]interface{}
