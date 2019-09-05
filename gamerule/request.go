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
		messagehandle.ErrorLogPrintln("game7", "SetInfo Index Error")
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
		"normalreel": normalScroll,
		"respinreel": {GetRespinScroll(RTPSetting)},
	}
	return scrollmap
}

// GetInitBetRate init info
func GetInitBetRate() interface{} {
	tmp := make(map[string]interface{})
	tmp["betrate"] = betRate
	tmp["betratelinkindex"] = betRateLinkIndex
	tmp["betratedefaultindex"] = betRateDefaultIndex
	tmp["winratearray"] = resultRateArray
	return tmp
}

// GetBetMoney ...
func GetBetMoney(betIndex int64) int64 {
	betrate := BetRate()
	return betrate[betIndex]
}

// BetRate ...
func BetRate() []int64 {
	return betRate
}

// GameRequest game server api return game result, game attach, totalwin
func GameRequest(playerID, betIndex int64, attach []map[string]interface{}) (map[string]interface{}, []map[string]interface{}, map[string]int64) {
	attinfo := attachDataToAttachInfo(playerID, attach)
	betMoney := GetBetMoney(betIndex)

	JackPartBonusPoolx2 := attinfo.JackPartBonusPoolx2
	JackPartBonusPoolx3 := attinfo.JackPartBonusPoolx3
	JackPartBonusPoolx5 := attinfo.JackPartBonusPoolx5

	result := logicResult(betMoney, &attinfo)

	attinfo.JackPartBonusPoolx2 += int64(float64(betMoney) * jackPortTex[2])
	attinfo.JackPartBonusPoolx3 += int64(float64(betMoney) * jackPortTex[1])
	attinfo.JackPartBonusPoolx5 += int64(float64(betMoney) * jackPortTex[0])

	otherdata := make(map[string]int64)
	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betMoney"] = betMoney
	otherdata["JackPartBonusx2"] = JackPartBonusPoolx2 - attinfo.JackPartBonusPoolx2
	otherdata["JackPartBonusx3"] = JackPartBonusPoolx3 - attinfo.JackPartBonusPoolx3
	otherdata["JackPartBonusx5"] = JackPartBonusPoolx5 - attinfo.JackPartBonusPoolx5

	return result, attachInfoToAttachData(attinfo), otherdata
}

// GetRespinScroll return respin scroll
func GetRespinScroll(ScrollIndex int) []int {
	if ScrollIndex == 2 {
		return respinScroll2
	}
	return respinScroll1

}
