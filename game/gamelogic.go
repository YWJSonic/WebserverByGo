package game

import (
	"encoding/json"
	"fmt"

	"gitlab.com/WeberverByGo/code"
	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	gameRules "gitlab.com/WeberverByGo/game/game5"
	"gitlab.com/WeberverByGo/mycache"
)

func outputGame(freecount int) (map[string]interface{}, int64) {
	var result map[string]interface{}
	var totalScores int64

	ScrollIndex, plate, otherdata := gameRules.NoprmalPlate()
	gr := gameRules.GameResult(plate)
	fmt.Println(ScrollIndex, plate, gr, otherdata)

	result["plateindex"] = ScrollIndex
	result["plate"] = plate
	result["freecount"] = freecount + otherdata[1]
	result["isrespin"] = otherdata[0] == 1
	result["isfreegame"] = foundation.InterfaceToInt(result["freecount"]) > gameRules.FreeGameTrigger

	result["scores"] = 0
	if len(gr) > 0 {
		result["scores"] = gr[0][2]
		totalScores += foundation.InterfaceToInt64(result["scores"])
	}

	return result, int64(gr[0][2])
}

func outputFreeSpin() ([]interface{}, int64) {
	var result []interface{}
	var totalScores int64

	for i, max := 0, 5; i < max; i++ {
		var freeresult map[string]interface{}
		ScrollIndex, plate, otherdata := gameRules.FreePlate()
		gr := gameRules.GameResult(plate)
		fmt.Println(ScrollIndex, plate, gr, otherdata)

		freeresult["plateindex"] = ScrollIndex
		freeresult["plate"] = plate
		freeresult["scores"] = 0
		if len(gr) > 0 {
			freeresult["scores"] = gr[0][2]
			totalScores += foundation.InterfaceToInt64(freeresult["scores"])
		}

		result = append(result, freeresult)
	}
	return result, totalScores
}

func outRespin() (map[string]interface{}, int64) {
	var result map[string]interface{}
	var totalscores int64

	ScrollIndex, plate, otherdata := gameRules.RespinPlate()
	gr := gameRules.RespinResult(plate)
	fmt.Println(ScrollIndex, plate, gr, otherdata)

	result["plateindex"] = ScrollIndex
	result["plate"] = plate
	result["scores"] = 0
	if len(gr) > 0 {
		result["scores"] = gr[0][1]
		totalscores += foundation.InterfaceToInt64(result["scores"])
	}

	return result, totalscores
}

// 0:free game count
func getGameInfo(playerID int64, gameIndex int64) gameInfo {
	var info gameInfo
	gameinfo := mycache.GetGameInfo(playerID)

	if gameinfo == nil {
		row, err := db.GetAttachKind(playerID, gameIndex)

		if len(row) > 0 && err.ErrorCode == code.OK {
			info = toGameInfo(row)
		} else {

			db.NewAttach(playerID, gameIndex, 0, 0)
			info = newGameInfo()
		}
	} else {
		if errMsg := json.Unmarshal(gameinfo.([]byte), &info); errMsg != nil {
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
			info.FreeCount = foundation.InterfaceToInt64(row["IValue"])
		}
	}
	return info
}
