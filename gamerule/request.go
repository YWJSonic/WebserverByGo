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

	// RespinSetting = foundation.InterfaceToInt(att["RespinSetting"])
}

// GetInitScroll ...
func GetInitScroll() interface{} {
	scrollmap := map[string][][]int{
		"normalreel": normalScroll,
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

	result := logicResult(betMoney, &attinfo)
	otherdata := make(map[string]int64)
	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betMoney"] = betMoney

	return result, attachInfoToAttachData(attinfo), otherdata
}
