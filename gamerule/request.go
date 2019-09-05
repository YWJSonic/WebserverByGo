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
		"normalreel":  getNormalScorll(),
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

	otherdata["totalwinscore"] = foundation.InterfaceToInt64(result["totalwinscore"])
	otherdata["betmoney"] = betMoney

	return result, attachInfoToAttachData(attinfo), otherdata
}

// ScotterGameRequest game server api return game result, game attach, totalwin
func ScotterGameRequest(playerID, betMoney, luckydrawselect, scotterID int64, attach []map[string]interface{}) (map[string]interface{}, []map[string]interface{}, map[string]interface{}) {
	attinfo := attachDataToAttachInfo(playerID, attach)
	scotterIndo := attinfo.ScotterInfos[scotterID]
	scotterIndo.DayScotterGameInfo = 1
	attinfo.ScotterInfos[scotterID] = scotterIndo
	var scotterCombination []int

	if luckydrawselect == 6 {
		scotterCombination = scotterGameMysteryIndexCombination[foundation.RangeRandom(scotterGameMysteryWeightings)]
	} else {
		scotterCombination = scotterGameMysteryIndexCombination[scotterGameDefaultCombinationIndex[luckydrawselect]]
	}

	result, otherdata := logicScotterGameResult(betMoney, scotterCombination[0], scotterCombination[1], &attinfo)
	result["scottercombination"] = scotterCombination
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

func getNormalScorll() [][]int {
	if RTPSetting == 6 {
		return normalScroll[5]
	} else if RTPSetting == 5 {
		return normalScroll[4]
	} else if RTPSetting == 4 {
		return normalScroll[3]
	} else if RTPSetting == 3 {
		return normalScroll[2]
	} else if RTPSetting == 2 {
		return normalScroll[1]
	} else {
		return normalScroll[0]
	}
}

func getNormalWildWinRateWeightings() []int64 {
	if RTPSetting == 6 {
		return normalWildWinRateWeightings[5]
	} else if RTPSetting == 5 {
		return normalWildWinRateWeightings[4]
	} else if RTPSetting == 4 {
		return normalWildWinRateWeightings[3]
	} else if RTPSetting == 3 {
		return normalWildWinRateWeightings[2]
	} else if RTPSetting == 2 {
		return normalWildWinRateWeightings[1]
	} else {
		return normalWildWinRateWeightings[0]
	}
}
