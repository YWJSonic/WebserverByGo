package log

import (
	"gitlab.com/ServerUtility/loginfo"
	"gitlab.com/WeberverByGo/db"
)

// SaveLog Save to db
func SaveLog(log loginfo.LogInfo) {
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
