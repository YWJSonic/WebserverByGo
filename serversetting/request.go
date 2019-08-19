package serversetting

import (
	"fmt"
	"sync"

	"gitlab.com/ServerUtility/cron"
	"gitlab.com/ServerUtility/foundation"
)

func init() {
	Setting = make(map[string]interface{})
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
	return serverTotalPayScore
}

// SetServerTotalPayScore ...
func SetServerTotalPayScore(value int64) {
	mu.Lock()
	defer mu.Unlock()
	serverTotalPayScore = value
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

// New ...
func New() map[string]interface{} {

	setting := map[string]interface{}{
		"servertime":   serverTime(),
		"maintaintime": maintainTime(),
	}
	return setting
}
