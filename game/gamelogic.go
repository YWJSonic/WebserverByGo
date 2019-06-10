package game

import (
	"encoding/json"
	"strings"

	"gitlab.com/WeberverByGo/log"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"

	"gitlab.com/WeberverByGo/foundation"
	gameRules "gitlab.com/WeberverByGo/game/game5"
	"gitlab.com/WeberverByGo/mycache"
)

// GetInitScroll ...
func GetInitScroll() interface{} {
	return gameRules.Scroll()
}

// GetBetRate ...
func GetBetRate() interface{} {
	return gameRules.BetRate()
}

// GetBetMoney ...
func GetBetMoney(betIndex int64) int64 {
	betrate := gameRules.BetRate()
	return betrate[betIndex]
}

// GameRequest ...
func gameRequest(playerID, betIndex int64) (map[string]interface{}, int64) {
	attach := GetAttach(playerID)
	betMoney := GetBetMoney(betIndex)
	result := gameRules.Result(betMoney, attach.FreeCount)

	attach.FreeCount = foundation.InterfaceToInt(result["freecount"])
	saveAttach(playerID, attach)

	msg := foundation.JSONToString(result)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	loginfo := log.New(log.GameResult)
	loginfo.PlayerID = playerID
	loginfo.IValue1 = foundation.InterfaceToInt64(result["totalwinscore"])
	loginfo.IValue2 = betMoney
	loginfo.Msg = msg
	log.SaveLog(loginfo)
	return result, foundation.InterfaceToInt64(result["totalwinscore"])
}

// InitAttach game start init attach
func InitAttach(playerid int64) {
	att := newAttach()
	saveAttach(playerid, att)
}

// GetAttach 0:free game count
func GetAttach(playerID int64) Attach {
	var info Attach
	gameIndex := gameRules.GameIndex
	attach := mycache.GetAttach(playerID)

	if attach == nil {
		// attach in db
		// row, err := db.GetAttachKind(playerID, gameIndex)
		// if len(row) > 0 && err.ErrorCode == code.OK {
		// 	// db data
		// 	info = toAttach(row)
		// } else {
		// 	// new data
		// 	db.NewAttach(playerID, gameIndex, 0, 0)
		// 	info = newAttach()
		// }

		info = newAttach()
	} else {
		// cache data
		if errMsg := json.Unmarshal(attach.([]byte), &info); errMsg != nil {
			errorlog.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(attach.([]byte)))
			info = newAttach()
		}
	}

	return info
}
func saveAttach(playerid int64, info Attach) {
	mycache.SetAttach(playerid, info.ToJSONStr())
	// db.UpdateAttach(playerid, gameRules.GameIndex, 0, info.FreeCount)
}
func newAttach() Attach {
	return Attach{}
}
func toAttach(rows []map[string]interface{}) Attach {
	var info Attach

	for _, row := range rows {
		switch foundation.InterfaceToInt64(row["Type"]) {
		case 0:
			info.FreeCount = foundation.InterfaceToInt(row["IValue"])
		}
	}
	return info
}
