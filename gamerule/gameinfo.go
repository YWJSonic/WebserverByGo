package gamerule

import "gitlab.com/ServerUtility/foundation"

// DayScotterGameCountKey DayScotterGameCount Key
const DayScotterGameCountKey = 1

// AttachInfo game att data
type AttachInfo struct {
	PlayerID             int64
	Kind                 int64
	DayScotterGameCount  int64
	FreeGameBetLockIndex map[int64]int64 //Key: DayScotterGameCount * 100 + 1 * 10 + WeekDay
	DayScotterGameInfo   map[int64]int64 //Key: DayScotterGameCount * 100 + 2 * 10 + WeekDay
	ScotterID            int64
}

func newAttchInfo() AttachInfo {
	return AttachInfo{
		DayScotterGameInfo:   make(map[int64]int64),
		FreeGameBetLockIndex: make(map[int64]int64),
	}
}

func attachDataToAttachInfo(playerID int64, att []map[string]interface{}) AttachInfo {
	var attType int64
	var attIValue int64
	attachInfo := newAttchInfo()
	// attachInfo.JackPartBonusPool = make(map[int64]int64)
	if len(att) > 0 {
		for _, row := range att {
			attIValue = foundation.InterfaceToInt64(row["IValue"])
			attType = foundation.InterfaceToInt64(row["Type"])
			attachInfo.PlayerID = foundation.InterfaceToInt64(row["PlayerID"])
			attachInfo.Kind = foundation.InterfaceToInt64(row["Kind"])

			if attType == 1 {
				attachInfo.DayScotterGameCount = attIValue
			} else if (attType%100)/10 == 1 {
				attachInfo.DayScotterGameInfo[attType] = attIValue
			} else if (attType%100)/10 == 2 {
				attachInfo.FreeGameBetLockIndex[attType] = attIValue
			}
		}
	} else {
		attachInfo.Kind = GameIndex
		attachInfo.PlayerID = playerID
	}
	return attachInfo
}

func attachInfoToAttachData(attachInfo AttachInfo) []map[string]interface{} {
	var att []map[string]interface{}
	playerID := attachInfo.PlayerID
	attKind := attachInfo.Kind

	att = append(att, toAttMap(playerID, attKind, DayScotterGameCountKey, attachInfo.DayScotterGameCount))
	for attType, attIValue := range attachInfo.FreeGameBetLockIndex {
		att = append(att, toAttMap(playerID, attKind, attType, attIValue))
	}
	for attType, attIValue := range attachInfo.DayScotterGameInfo {
		att = append(att, toAttMap(playerID, attKind, attType, attIValue))
	}

	return att
}

func toAttMap(playerid, attKind, attType, iValue int64) map[string]interface{} {
	return map[string]interface{}{
		"PlayerID": playerid,
		"Kind":     attKind,
		"Type":     attType,
		"IValue":   iValue,
	}
}

// NewDayScotterGameID ...
func NewDayScotterGameID(DayScotterGameCount int64) int64 {
	return int64(foundation.ServerNow().Weekday()) + (DayScotterGameCount * 10)
}

// DayScotterGameInfoKey ...
func DayScotterGameInfoKey(scotterID int64) int64 {
	return scotterID + 10
}

// FreeGameBetLockIndexKey ...
func FreeGameBetLockIndexKey(scotterID int64) int64 {
	return scotterID + 20
}
