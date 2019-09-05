package gamerule

import (
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/messagehandle"
)

// InitAttach init game attach
func InitAttach(playerID int64) {

}

// ConvertToGameAttach ...
func ConvertToGameAttach(playerID int64, attinfo []map[string]interface{}) AttachInfo {
	return attachDataToAttachInfo(playerID, attinfo)
}

// SetInfo ...
func SetInfo(gameIndex int, att map[string]interface{}) {
	if GameIndex != gameIndex {
		messagehandle.ErrorLogPrintln("SetInfo Error:", gameIndex, att)
		return
	}

	if value, ok := att["RespinSetting"]; ok {
		RespinSetting = foundation.InterfaceToInt(value)
	}
	if value, ok := att["RTPSetting"]; ok {
		RTPSetting = foundation.InterfaceToInt(value)
	}
	if value, ok := att["WinScoreLimit"]; ok {
		WinScoreLimit = foundation.InterfaceToInt64(value)
	}
	if value, ok := att["WinBetRateLimit"]; ok {
		WinBetRateLimit = foundation.InterfaceToInt64(value)
	}
}

// GetInitScroll ...
func GetInitScroll() interface{} {
	scrollmap := map[string][][]int{
		"normalreel": getNormalScorll(),
		"freereel":   freeScroll,
	}
	return scrollmap
}

// GetInitBetRate init info
func GetInitBetRate() interface{} {
	tmp := make(map[string]interface{})
	tmp["betrate"] = betRate
	tmp["betratelinkindex"] = betRateLinkIndex
	tmp["betratedefaultindex"] = betRateDefaultIndex
	return tmp
}

// GetBetMoney ...
func GetBetMoney(betIndex int64) int64 {
	betrate := BetRate()
	if betIndex >= int64(len(betrate)) {
		return betrate[betRateDefaultIndex]
	}
	return betrate[betIndex]
}

// BetRate ...
func BetRate() []int64 {
	return betRate
}

// GameRequest game server api return game result, game attach, totalwin
func GameRequest(playerID, betIndex int64, attach []map[string]interface{}) (map[string]interface{}, []map[string]interface{}, map[string]interface{}) {
	attinfo := attachDataToAttachInfo(playerID, attach)
	betMoney := GetBetMoney(betIndex)

	result, otherdata := logicResult(betMoney, &attinfo)

	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betmoney"] = betMoney

	return result, attachInfoToAttachData(attinfo), otherdata
}

func getNormalScorll() [][]int {
	if RTPSetting == 6 {
		return normalScroll[RTPSetting-1]
	} else if RTPSetting == 5 {
		return normalScroll[RTPSetting-1]
	} else if RTPSetting == 4 {
		return normalScroll[RTPSetting-1]
	} else if RTPSetting == 3 {
		return normalScroll[RTPSetting-1]
	} else if RTPSetting == 2 {
		return normalScroll[RTPSetting-1]
	} else {
		return normalScroll[0]
	}
}
