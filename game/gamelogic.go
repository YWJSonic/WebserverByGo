package game

import (
	"encoding/json"
	"strings"

	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/messagehandle/log"

	"gitlab.com/WeberverByGo/foundation"
	gameRules "gitlab.com/WeberverByGo/game/game5"
	"gitlab.com/WeberverByGo/mycache"
)

// GetInitScroll ...
func GetInitScroll() interface{} {
	return gameRules.Scroll()
}

// GetInitBetRate ...
func GetInitBetRate() interface{} {
	return gameRules.BetInitInfo()
}

// GetBetMoney ...
func GetBetMoney(betIndex int) int64 {
	betrate := gameRules.BetRate()
	return betrate[betIndex]
}

// GameRequest ...
func gameRequest(playerID int64, betIndex int) (map[string]interface{}, int64) {
	attach := GetAttach(playerID)
	if attach.FreeCount >= gameRules.FreeGameTrigger {
		attach.FreeCount %= gameRules.FreeGameTrigger
		attach = newAttach()
	}
	if attach.IsLockBet != 0 {
		betIndex = attach.LockBetIndex
	}

	betMoney := GetBetMoney(betIndex)
	result := gameRules.Result(betMoney, attach.FreeCount)
	freeCount := foundation.InterfaceToInt(result["freecount"])
	attach.FreeCount = freeCount

	result["islockbet"] = 0
	result["lockbetindex"] = 1
	if freeCount%gameRules.FreeGameTrigger > 0 {
		attach = lockBet(attach, betIndex)
		result["islockbet"] = 1
		result["lockbetindex"] = betIndex
	} else if freeCount%gameRules.FreeGameTrigger == 0 && freeCount > 0 {
		attach = unlockBet(attach)
	}

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

func lockBet(att Attach, betindex int) Attach {
	att.LockBetIndex = betindex
	att.IsLockBet = 1
	return att
}
func unlockBet(att Attach) Attach {
	att.LockBetIndex = 1
	att.IsLockBet = 0
	return att
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
	mycache.SetAttach(playerid, foundation.JSONToString(info))
	// db.UpdateAttach(playerid, gameRules.GameIndex, 0, info.FreeCount)
}
func newAttach() Attach {
	return Attach{
		FreeCount:    0,
		IsLockBet:    0,
		LockBetIndex: 1,
	}
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
