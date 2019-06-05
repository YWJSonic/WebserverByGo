package game

import (
	"encoding/json"

	"gitlab.com/WeberverByGo/messagehandle/errorlog"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/db"
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

// GameRequest ...
func gameRequest(playerID, betIndex int64) (map[string]interface{}, int64) {
	attach := GetAttach(playerID)
	result := gameRules.Result(betIndex, attach.FreeCount)

	attach.FreeCount = foundation.InterfaceToInt(result.(map[string]interface{})["freecount"])
	if attach.FreeCount >= gameRules.FreeGameTrigger {
		attach.FreeCount %= gameRules.FreeGameTrigger
	}
	saveAttach(1, gameRules.GameIndex(), attach)
	return result.(map[string]interface{}), foundation.InterfaceToInt64(result.(map[string]interface{})["totalwinscore"])
}

// GetAttach 0:free game count
func GetAttach(playerID int64) Attach {
	var info Attach
	gameIndex := gameRules.GameIndex()
	attach := mycache.GetAttach(playerID)

	if attach == nil {
		row, err := db.GetAttachKind(playerID, gameIndex)

		if len(row) > 0 && err.ErrorCode == code.OK {
			// db data
			info = toAttach(row)
		} else {

			// new data
			db.NewAttach(playerID, gameIndex, 0, 0)
			info = newAttach()
		}
	} else {
		// cache data
		if errMsg := json.Unmarshal(attach.([]byte), &info); errMsg != nil {
			errorlog.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(attach.([]byte)))
			info = newAttach()
		}
	}

	return info
}
func saveAttach(playerid int64, gameIndex int64, info Attach) {
	mycache.SetAttach(playerid, info.ToJSONStr())
	db.UpdateAttach(playerid, gameIndex, 0, info.FreeCount)
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
