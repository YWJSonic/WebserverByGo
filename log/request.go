package log

import (
	"time"

	"gitlab.com/WeberverByGo/db"
)

// New Default log
func New(ActivityEvent int) LogInfo {
	return LogInfo{
		Time:          time.Now().Unix(),
		ActivityEvent: ActivityEvent,
	}
}

// SaveLog Save to db
func SaveLog(log LogInfo) {
	db.SetLog(
		log.Account,
		log.PlayerID,
		log.Time,
		log.ActivityEvent,
		log.IValue1,
		log.IValue2,
		log.IValue3,
		log.SValue1,
		log.SValue2,
		log.SValue3,
		log.Msg)
}
