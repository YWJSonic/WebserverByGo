package log

import (
	"time"

	"gitlab.com/WeberverByGoGame5/db"
)

// New Default log
func New(activityEvent int) LogInfo {
	return LogInfo{
		Time:          time.Now().Unix(),
		ActivityEvent: activityEvent,
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
