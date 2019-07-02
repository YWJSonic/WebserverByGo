package log

import (
	"gitlab.com/ServerUtility/loginfo"
	db "gitlab.com/WeberverByGo/handledb"
)

// AcouuntLogin log account
func AcouuntLogin(GameAccount string) {

	loginfo := loginfo.New(loginfo.Login)
	loginfo.Account = GameAccount
	SaveLog(loginfo)

}

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
