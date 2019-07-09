package gamerule

import "gitlab.com/ServerUtility/foundation"

// AttachInfo game att data
type AttachInfo struct {
	PlayerID            int64
	Kind                int64
	JackPartBonusPoolx2 int64
	JackPartBonusPoolx3 int64
	JackPartBonusPoolx5 int64
}

func attachDataToAttachInfo(playerID int64, att []map[string]interface{}) AttachInfo {
	var attachInfo AttachInfo
	var attType int64
	var attIValue int64
	// attachInfo.JackPartBonusPool = make(map[int64]int64)
	if len(att) > 0 {
		for _, row := range att {
			attType = foundation.InterfaceToInt64(row["Type"])
			attIValue = foundation.InterfaceToInt64(row["IValue"])
			attachInfo.PlayerID = foundation.InterfaceToInt64(row["PlayerID"])
			attachInfo.Kind = foundation.InterfaceToInt64(row["Kind"])

			if attType == 0 {
				attachInfo.JackPartBonusPoolx2 = attIValue
			} else if attType == 1 {
				attachInfo.JackPartBonusPoolx3 = attIValue
			} else if attType == 2 {
				attachInfo.JackPartBonusPoolx5 = attIValue
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

	att = append(att, toAttMap(playerID, attKind, 0, attachInfo.JackPartBonusPoolx2))
	att = append(att, toAttMap(playerID, attKind, 1, attachInfo.JackPartBonusPoolx3))
	att = append(att, toAttMap(playerID, attKind, 2, attachInfo.JackPartBonusPoolx5))

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
