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
func gameRequest(BetIndex int64) map[string]interface{} {

	gameinfo := getGameInfo(1, gameRules.GameIndex())
	result := gameRules.Result(BetIndex, gameinfo.FreeCount)

	normalresult := result.(map[string]interface{})["normalresult"]
	gameinfo.FreeCount = foundation.InterfaceToInt(normalresult.(map[string]interface{})["freecount"])
	saveGameInfo(1, gameRules.GameIndex(), gameinfo)
	return result.(map[string]interface{})
}

// 0:free game count
func getGameInfo(playerID int64, gameIndex int64) gameInfo {
	var info gameInfo
	gameinfo := mycache.GetGameInfo(playerID)

	if gameinfo == nil {
		row, err := db.GetAttachKind(playerID, gameIndex)

		if len(row) > 0 && err.ErrorCode == code.OK {
			// db data
			info = toGameInfo(row)
		} else {

			// new data
			db.NewAttach(playerID, gameIndex, 0, 0)
			info = newGameInfo()
		}
	} else {
		// cache data
		if errMsg := json.Unmarshal(gameinfo.([]byte), &info); errMsg != nil {
			errorlog.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(gameinfo.([]byte)))
			info = newGameInfo()
		}
	}

	return info
}
func saveGameInfo(playerid int64, gameIndex int64, info gameInfo) {
	mycache.SetGameInfo(playerid, info.ToJSONStr())
	db.UpdateAttach(playerid, gameIndex, 0, info.FreeCount)
}

func newGameInfo() gameInfo {
	return gameInfo{}
}
func toGameInfo(rows []map[string]interface{}) gameInfo {
	var info gameInfo

	for _, row := range rows {
		switch foundation.InterfaceToInt64(row["Type"]) {
		case 0:
			info.FreeCount = foundation.InterfaceToInt(row["IValue"])
		}
	}
	return info
}
