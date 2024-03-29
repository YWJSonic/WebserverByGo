package attach

import (
	"encoding/json"

	"gitlab.com/ServerUtility/code"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
	mycache "gitlab.com/WebserverByGoBase/handlecache"
	db "gitlab.com/WebserverByGoBase/handledb"
)

// GetAttach 0:free game count
func GetAttach(playerID int64, gameIndex int64, isSaveToDB bool) []map[string]interface{} {
	var info []map[string]interface{}
	attach := mycache.GetAttach(playerID)
	if attach == nil {
		if isSaveToDB {
			info = loadAttachFromDB(playerID, gameIndex)
		}

		if info == nil {
			info = newAttach()
		}

	} else {
		// cache data
		if errMsg := json.Unmarshal(attach.([]byte), &info); errMsg != nil {
			messagehandle.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(attach.([]byte)))
			info = newAttach()
		}
	}

	return info
}

// GetAttachByType ...
func GetAttachByType(playerID int64, gameIndex int64, attType int64, isSaveToDB bool) []map[string]interface{} {
	var info []map[string]interface{}
	attach := mycache.GetAttach(playerID)
	if attach == nil {
		if isSaveToDB {
			info = loadAttachByTypeFromDB(playerID, gameIndex, attType)
		}

		if info == nil {
			info = newAttach()
		}

	} else {
		// cache data
		if errMsg := json.Unmarshal(attach.([]byte), &info); errMsg != nil {
			messagehandle.ErrorLogPrintln("GameLogic", playerID, gameIndex, string(attach.([]byte)))
			info = newAttach()
		}
	}

	return info
}

// GetAttachByTypeRange ...
func GetAttachByTypeRange(playerID, gameIndex, miniAttType, maxAttachType int64) []map[string]interface{} {
	return loadAttachByTypeRangeFromDB(playerID, gameIndex, miniAttType, maxAttachType)
}

// SaveAttach ...
func SaveAttach(playerid int64, gameIndex int, info []map[string]interface{}, isSaveToDB bool) {
	saveToCache(playerid, gameIndex, info)
	if isSaveToDB {
		SaveAttachToDB(playerid, gameIndex, info)
	}
}

// SaveAttachToDB ...
func SaveAttachToDB(playerid int64, gameIndex int, info []map[string]interface{}) {
	for _, value := range info {
		saveAttachToDB(playerid, gameIndex, foundation.InterfaceToInt64(value["Type"]), foundation.InterfaceToInt64(value["IValue"]))
	}
}

func loadAttachByTypeRangeFromDB(playerid, gameIndex, miniAttType, maxAttachType int64) []map[string]interface{} {
	row, err := db.GetAttachTypeRange(playerid, gameIndex, miniAttType, maxAttachType)
	if err.ErrorCode != code.OK {
		return nil
	}
	return row
}

func loadAttachByTypeFromDB(playerid, gameIndex, attType int64) []map[string]interface{} {
	row, err := db.GetAttachType(playerid, gameIndex, attType)
	if err.ErrorCode != code.OK {
		return nil
	}
	return row
}

func loadAttachFromDB(playerid, gameIndex int64) []map[string]interface{} {
	row, err := db.GetAttachKind(playerid, gameIndex)
	if err.ErrorCode != code.OK {
		return nil
	}
	return row
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
