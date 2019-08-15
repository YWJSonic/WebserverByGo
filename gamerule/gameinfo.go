package gamerule

import "gitlab.com/ServerUtility/foundation"

// DayScotterGameCountKey DayScotterGameCount Key
const DayScotterGameCountKey = 1

// AttachInfo game att data
type AttachInfo struct {
	PlayerID             int64
	Kind                 int64
	DayScotterGameCount  int64
	FreeGameBetLockMoney map[int64]int64 //Key: DayScotterGameCount * 100 + 1 * 10 + WeekDay
	DayScotterGameInfo   map[int64]int64 //Key: DayScotterGameCount * 100 + 2 * 10 + WeekDay
	ScotterInfos         map[int64]ScotterInfo
	ScotterID            int64
}

// NewScotterInfo ...
func (info *AttachInfo) NewScotterInfo(scotterID, dayScotterGameInfo, freeGameBetLockMoney int64) {
	scotterInfo := ScotterInfo{
		ScotterID:            scotterID,
		DayScotterGameInfo:   dayScotterGameInfo,
		FreeGameBetLockMoney: freeGameBetLockMoney,
	}
	info.ScotterInfos[scotterID] = scotterInfo
}

// NewDayScotterID ...
func (info *AttachInfo) NewDayScotterID() int64 {
	info.DayScotterGameCount++
	return (info.DayScotterGameCount * 100) + int64(foundation.ServerNow().Weekday()*10)
}

// ScotterInfo ...
type ScotterInfo struct {
	ScotterID            int64
	FreeGameBetLockMoney int64
	DayScotterGameInfo   int64
}

func newAttchInfo() AttachInfo {
	return AttachInfo{
		DayScotterGameInfo:   make(map[int64]int64),
		FreeGameBetLockMoney: make(map[int64]int64),
		ScotterInfos:         make(map[int64]ScotterInfo),
	}
}

func attachDataToAttachInfo(playerID int64, att []map[string]interface{}) AttachInfo {
	var attType int64
	var attIValue int64
	attachInfo := newAttchInfo()

	if len(att) > 0 {
		for _, row := range att {
			attIValue = foundation.InterfaceToInt64(row["IValue"])
			attType = foundation.InterfaceToInt64(row["Type"])
			attachInfo.PlayerID = foundation.InterfaceToInt64(row["PlayerID"])
			attachInfo.Kind = foundation.InterfaceToInt64(row["Kind"])

			if attType == 1 {
				attachInfo.DayScotterGameCount = attIValue
			} else if (attType % 10) == 1 {
				scotterID := GameInfoKeyToScotterID(attType)

				if Info, ok := attachInfo.ScotterInfos[scotterID]; ok {
					Info.DayScotterGameInfo = attIValue
					attachInfo.ScotterInfos[scotterID] = Info
				} else {
					Info := ScotterInfo{
						DayScotterGameInfo: attIValue,
						ScotterID:          scotterID,
					}
					attachInfo.ScotterInfos[scotterID] = Info
				}
			} else if (attType % 10) == 2 {
				scotterID := LockIndexKeyToScotterID(attType)

				if Info, ok := attachInfo.ScotterInfos[scotterID]; ok {
					Info.FreeGameBetLockMoney = attIValue
					attachInfo.ScotterInfos[scotterID] = Info
				} else {
					Info := ScotterInfo{
						FreeGameBetLockMoney: attIValue,
						ScotterID:            scotterID,
					}
					attachInfo.ScotterInfos[scotterID] = Info

				}
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

	// var tmpLockID int64
	// for attType, attIValue := range attachInfo.DayScotterGameInfo {
	// 	att = append(att, toAttMap(playerID, attKind, attType, attIValue))
	// 	tmpLockID = FreeGameBetLockIndexKey(GameInfoKeyToScotterID(attType))
	// 	att = append(att, toAttMap(playerID, attKind, tmpLockID, attachInfo.FreeGameBetLockMoney[tmpLockID]))
	// }

	for scotterID, Info := range attachInfo.ScotterInfos {
		att = append(att, toAttMap(playerID, attKind, DayScotterGameInfoKey(scotterID), Info.DayScotterGameInfo))
		att = append(att, toAttMap(playerID, attKind, FreeGameBetLockIndexKey(scotterID), Info.FreeGameBetLockMoney))
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

// DayScotterGameInfoKey ...
func DayScotterGameInfoKey(scotterID int64) int64 {
	return scotterID + 1
}

// GameInfoKeyToScotterID ...
func GameInfoKeyToScotterID(GameInfoKey int64) int64 {
	return GameInfoKey - 1
}

// FreeGameBetLockIndexKey ...
func FreeGameBetLockIndexKey(scotterID int64) int64 {
	return scotterID + 2
}

// LockIndexKeyToScotterID ...
func LockIndexKeyToScotterID(GameInfoKey int64) int64 {
	return GameInfoKey - 2
}
