package serversetting

import (
	"fmt"
	"sync"
	"time"

	cron "gitlab.com/ServerUtility/cron.v3"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/settinginfo"
	db "gitlab.com/WeberverByGoGame8/handledb"
)

var mu *sync.RWMutex

// CacheDeleteTime cache keep time
const CacheDeleteTime time.Duration = time.Hour

// ConnectTimeOut Client connect time out
const ConnectTimeOut int64 = 15e9

// platform api url
var (
	// MaintainStartTime cron maintain schedule
	MaintainStartTime = "0 15 * * *"
	// MaintainFinishTime cron maintain schedule
	MaintainFinishTime = "0 16 * * *"
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
var (
	// ServerTotalPayScore All player win score
	serverTotalPayScore settinginfo.Info
)

// Maintain Is sow maintain time
var maintain = false

// Setting settint from db data
// var Setting map[string]settinginfo.Info

func init() {
	mu = new(sync.RWMutex)
}

// ServerURL ...
func ServerURL() string {
	return fmt.Sprintf("%s:%s", IP, PORT)
}

// IsMaintain ...
func IsMaintain() bool {
	mu.RLock()
	defer mu.RUnlock()
	return maintain
}

// EnableMaintain ...
func EnableMaintain(enable bool) {
	mu.Lock()
	defer mu.Unlock()
	maintain = enable
}

// GetServerTotalPayScore ...
func GetServerTotalPayScore() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return serverTotalPayScore.IValue
}

// SetServerTotalPayScore ...
func SetServerTotalPayScore(value int64) {
	mu.Lock()
	defer mu.Unlock()
	serverTotalPayScore.IValue = value
}

// ServerTime ...
func serverTime() int64 {
	return foundation.ServerNowTime()
}

// MaintainTime ...
func maintainTime() int64 {
	target, _ := cron.ParseStandard(MaintainStartTime)
	return target.Next(foundation.ServerNow()).Unix()
}

// InsertDBSetting update setting by db setting data
func InsertDBSetting(dbDatas []map[string]interface{}, GameIndex int) {
	for _, settingData := range dbDatas {
		info := settinginfo.ConvertToInfo(settingData)
		if info.Key == foundation.ServerTotalPayScoreKey(GameIndex) {
			serverTotalPayScore = info
		}
	}
}

// RefreshDBSetting ...
func RefreshDBSetting(gameTypeIndex int, serverDayPayDefault int64) {
	timeslip := foundation.ServerNowTime() - serverTotalPayScore.LastRefulsh
	timelimit := int64(time.Hour * 24 / time.Second)

	if timeslip >= timelimit {
		db.ReflushSetting(foundation.ServerTotalPayScoreKey(gameTypeIndex), serverDayPayDefault, "")
	}
}

// New ...
func New() map[string]interface{} {
	return map[string]interface{}{
		"servertime":   serverTime(),
		"maintaintime": maintainTime(),
	}
}
