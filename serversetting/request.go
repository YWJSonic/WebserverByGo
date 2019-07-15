package serversetting

import (
	"fmt"
	"sync"

	"github.com/robfig/cron"
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
