package attach

import (
	"encoding/json"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	mycache "gitlab.com/WeberverByGo/handlecache"
	db "gitlab.com/WeberverByGo/handledb"
)

// GetAttach 0:free game count
func GetAttach(playerID int64, gameIndex int) []map[string]interface{} {
	var info []map[string]interface{}
	attach := mycache.GetAttach(playerID)
	if attach == nil {
		info = newAttach()
	} else {
		// cache data
		if errMsg := json.Unmarshal(attach.([]byte), &info); errMsg != nil {
			messagehandle.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(attach.([]byte)))
			info = newAttach()
		}
	}

	return info
}

// SaveAttach ...
func SaveAttach(playerid int64, gameIndex int, info []map[string]interface{}, isSaveToDB bool) {
	saveToCache(playerid, gameIndex, info)
	if isSaveToDB {
		for _, value := range info {
			saveAttachToDB(playerid, gameIndex, foundation.InterfaceToInt64(value["Type"]), foundation.InterfaceToInt64(value["IValue"]))
		}
	}
}

func saveToCache(playerid int64, gameIndex int, info []map[string]interface{}) {
	mycache.SetAttach(playerid, foundation.JSONToString(info))
}
func saveAttachToDB(playerid int64, gameIndex int, atttype, value int64) {
	db.UpdateAttach(playerid, gameIndex, atttype, value)
}
func newAttach() []map[string]interface{} {
	var attmap []map[string]interface{}
	return attmap
}
