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

	// RespinSetting = foundation.InterfaceToInt(att["RespinSetting"])
}

// GetInitScroll ...
func GetInitScroll() interface{} {
	scrollmap := map[string][][]int{
		"normalreel":  normalScroll,
		"scotterreel": scotterScroll,
	}
	return scrollmap
}

// GetInitBetRate init info
func GetInitBetRate() interface{} {
	tmp := make(map[string]interface{})
	tmp["betrate"] = betRate
	tmp["betratelinkindex"] = betRateLinkIndex
	tmp["betratedefaultindex"] = betRateDefaultIndex
	tmp["luckydrawcombination"] = luckydrawInit()
	tmp["luckydrawwinrate"] = scotterGameWildWinRate
	tmp["luckydrawspin"] = scotterGameSpinTime
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

	if value, ok := otherdata["isscotter"]; ok && value == 1 {
		attinfo.DayScotterGameCount++
		weekDay := int64(foundation.ServerNow().Weekday())
		scotterid := attinfo.DayScotterGameCount*100 + weekDay
		attinfo.DayScotterGameInfo[scotterid+10] = 0
		attinfo.FreeGameBetLockIndex[scotterid+20] = betMoney
		result["isscotter"] = 1
		result["scotterid"] = scotterid
	} else {
		result["isscotter"] = 0
		result["scotterid"] = 0
	}

	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betmoney"] = betMoney

	return result, attachInfoToAttachData(attinfo), otherdata
}

// ScotterGameRequest game server api return game result, game attach, totalwin
func ScotterGameRequest(playerID, betMoney, luckydrawselect int64, attach []map[string]interface{}) (map[string]interface{}, []map[string]interface{}, map[string]interface{}) {
	attinfo := attachDataToAttachInfo(playerID, attach)
	var scotterCombination []int

	if luckydrawselect == 6 {
		scotterCombination = scotterGameMysteryIndexCombination[foundation.RangeRandom(scotterGameMysteryWeightings)]
	} else {
		scotterCombination = scotterGameMysteryIndexCombination[scotterGameDefaultCombinationIndex[luckydrawselect]]
	}

	result, otherdata := logicScotterGameResult(betMoney, scotterCombination[0], scotterCombination[1], &attinfo)

	if value, ok := otherdata["isscotter"]; ok && value == 1 {
		attinfo.DayScotterGameCount++
		weekDay := int64(foundation.ServerNow().Weekday())
		scotterid := attinfo.DayScotterGameCount*100 + weekDay
		attinfo.DayScotterGameInfo[scotterid+10] = 0
		attinfo.FreeGameBetLockIndex[scotterid+20] = betMoney
		result["isscotter"] = 1
		result["scotterid"] = scotterid
	}

	// otherdata := make(map[string]int64)
	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betmoney"] = betMoney

	return result, attachInfoToAttachData(attinfo), otherdata
}

func luckydrawInit() [][]int {
	var result [][]int
	for _, DefaultCombination := range scotterGameDefaultCombinationIndex {
		result = append(result, scotterGameMysteryIndexCombination[DefaultCombination])
	}
	return result
}
