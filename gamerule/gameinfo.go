package gamerule

// DayScotterGameCountKey DayScotterGameCount Key
const DayScotterGameCountKey = 1

// AttachInfo game att data
type AttachInfo struct {
	PlayerID int64
	Kind     int64
}

func newAttchInfo() AttachInfo {
	return AttachInfo{}
}

func attachDataToAttachInfo(playerID int64, att []map[string]interface{}) AttachInfo {
	attachInfo := newAttchInfo()

	return attachInfo
}

func attachInfoToAttachData(attachInfo AttachInfo) []map[string]interface{} {
	var att []map[string]interface{}

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

// FreeGameBetLockIndexKey ...
func FreeGameBetLockIndexKey(scotterID int64) int64 {
	return scotterID + 2
}
