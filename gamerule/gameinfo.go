package gamerule

// AttachInfo game att data
type AttachInfo struct {
	PlayerID int64
	Kind     int64
}

func attachDataToAttachInfo(playerID int64, att []map[string]interface{}) AttachInfo {
	var attachInfo AttachInfo

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
