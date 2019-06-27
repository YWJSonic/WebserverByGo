package attach

import (
	"encoding/json"

	"gitlab.com/WeberverByGo/db"
	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/messagehandle/errorlog"
	"gitlab.com/WeberverByGo/mycache"
)

// GetAttach 0:free game count
func GetAttach(playerID int64, gameIndex int) map[string]interface{} {
	var info map[string]interface{}
	attach := mycache.GetAttach(playerID)
	if attach == nil {
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

// SaveAttach ...
func SaveAttach(playerid int64, gameIndex int, info map[string]interface{}, isSaveToDB bool) {
	saveAttach(playerid, gameIndex, info)
	if isSaveToDB {
		// for key, value := range info {
		// saveAttachToDB(playerid, gameIndex, info, key, value)
		// }
	}
}

func saveAttach(playerid int64, gameIndex int, info map[string]interface{}) {
	mycache.SetAttach(playerid, foundation.JSONToString(info))
}
func saveAttachToDB(playerid int64, gameIndex, atttype, int, value int64) {
	db.UpdateAttach(playerid, gameIndex, atttype, value)
}
func newAttach() map[string]interface{} {
	var attmap map[string]interface{}
	return attmap
}
