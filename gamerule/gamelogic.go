package gamerule

import (
	"fmt"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gameplate"
)

// Result att 0: freecount
func logicResult(betMoney int64, attinfo *AttachInfo) map[string]interface{} {
	var result = make(map[string]interface{})
	var totalWin int64

	normalresult, otherdata, normaltotalwin := outputGame(betMoney, attinfo)
	result = foundation.AppendMap(result, otherdata)
	result["normalresult"] = normalresult
	totalWin += normaltotalwin

	if otherdata["isrespin"].(int) == 1 {
		respinresult, respintotalwin := outRespin(betMoney, attinfo)
		totalWin += respintotalwin
		result["respin"] = respinresult
		result["isrespin"] = 1
	}

	result["totalwinscore"] = totalWin
	return result
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var WinRateIndex int
	var result map[string]interface{}
	otherdata := make(map[string]interface{})
	islink := false

	ScrollIndex, plate := gameplate.NewPlate(scrollSize, normalScroll)
	gameresult := normalResultArray(plate)

	otherdata["isrespin"] = 0

	if isRespin(plate) {
		otherdata["isrespin"] = 1
	}

	if len(gameresult) > 0 {
		islink = true
		WinRateIndex = gameresult[0][3]
		reGameResult := dynamicScore(plate, gameresult[0])
		switch WinRateIndex {
		case -101:
			totalScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx2
			attinfo.JackPartBonusPoolx2 = 0
		case -102:
			totalScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx3
			attinfo.JackPartBonusPoolx3 = 0
		case -103:
			totalScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx5
			attinfo.JackPartBonusPoolx5 = 0
		default:
			totalScores = betMoney * int64(reGameResult[3])
			switch plate[1] {
			case wild2:
				totalScores *= spWhildWinRate[0]
			case wild3:
				totalScores *= spWhildWinRate[1]
			case wild4:
				totalScores *= spWhildWinRate[2]
			default:
			}
		}
	}

	if totalScores < 0 {
		fmt.Println(totalScores)
	}
	result = gameplate.ResultMap(ScrollIndex, plate, totalScores, islink)
	return result, otherdata, totalScores
}

// outRespin out put respin result and totalwin
func outRespin(betMoney int64, attinfo *AttachInfo) ([]interface{}, int64) {
	var totalScores, respinScores, totalWinRate int64
	var WinRateIndex int
	var ScrollIndex, plate []int
	var result []interface{}
	respintScrollData := GetRespinScroll(RespinSetting)
	islink := false

	for index, max := 0, 200; index < max; index++ {
		islink = false
		respinScores = 0
		ScrollIndex, plate = gameplate.NewPlate([]int{1, 1, 1}, [][]int{{0}, respintScrollData, {0}})
		gameresult := respinResultArray(plate)

		if len(gameresult) > 0 {
			islink = true
			WinRateIndex = gameresult[0][3]
			reGameResult := dynamicScore(plate, gameresult[0])
			totalWinRate += int64(reGameResult[3])
			switch WinRateIndex {
			case -101:
				respinScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx2
				attinfo.JackPartBonusPoolx2 = 0
			case -102:
				respinScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx3
				attinfo.JackPartBonusPoolx3 = 0
			case -103:
				respinScores = betMoney*int64(reGameResult[3]) + attinfo.JackPartBonusPoolx5
				attinfo.JackPartBonusPoolx5 = 0
			default:
				respinScores = betMoney * int64(reGameResult[3])
			}
		}

		totalScores += respinScores
		freeresult := gameplate.ResultMap(ScrollIndex, plate, respinScores, islink)
		result = append(result, freeresult)

		if len(gameresult) <= 0 {
			break
		} else if index >= max-1 {
			result = append(result, emptyResult())
			break
		}
	}
	return result, totalScores
}

// winresultArray ...
func normalResultArray(plate []int) [][]int {
	var result [][]int

	for _, JackPortResult := range jackPortResults {
		if isJackportWin(plate, JackPortResult) {
			result = append(result, JackPortResult)
			if isSingleLine {
				return result
			}
		}
	}

	for _, ItemResult := range itemResults {
		if isNormalWin(plate, ItemResult) {
			result = append(result, ItemResult)
			if isSingleLine {
				return result
			}
		}
	}
	return result

}

// RespinResult result 0: icon index, 1: win rate
func respinResultArray(plate []int) [][]int {
	var result [][]int

	for _, JackPortResult := range jackPortResults {
		if isJackportWin(plate, JackPortResult) {
			result = append(result, JackPortResult)
			if isSingleLine {
				return result
			}
		}
	}

	for _, RespinResult := range respinitemResults {
		if isRespinWin(plate, RespinResult) {
			result = append(result, RespinResult)
			if isSingleLine {
				return result
			}
		}
	}

	return result
}

// EmptyResult return a not win result
func emptyResult() map[string]interface{} {
	return gameplate.ResultMap([]int{0, 0, 0}, []int{0, space, 0}, 0, false)
}

// dynamicScore convert results list dynamic score
func dynamicScore(plant, currendResult []int) []int {
	if !isDynamicResult(currendResult) {
		return currendResult
	}

	dynamicresult := make([]int, len(currendResult))
	copy(dynamicresult, currendResult)

	switch currendResult[3] {
	case -100:
		for _, result := range itemResults {
			if result[1] == plant[1] {
				dynamicresult[3] = result[3]
				break
			}
		}
	case -101:
		dynamicresult[3] = jackPartWinRate[0]
		break
	case -102:
		dynamicresult[3] = jackPartWinRate[1]
		break
	case -103:
		dynamicresult[3] = jackPartWinRate[2]
		break
	}

	return dynamicresult
}

func isDynamicResult(result []int) bool {
	if result[3] < 0 {
		return true
	}
	return false
}

func isNormalWin(plates []int, result []int) bool {
	IsWin := false
	for i, plate := range plates {
		IsWin = false

		if plate == space {
			return false
		}

		if plate == wild1 || plate == wild2 || plate == wild3 || plate == wild4 {
			IsWin = true
		} else {

			switch result[i] {
			case plate:
				IsWin = true
			case -1000:
				IsWin = true
			case -1001: // any bar
				if foundation.IsInclude(plate, symbolGroup[result[i]]) {
					IsWin = true
				}
			}
		}
		if !IsWin {
			return IsWin
		}
	}

	return IsWin
}

func isRespinWin(plates []int, result []int) bool {
	return isNormalWin(plates, result)
}

func isJackportWin(plates []int, result []int) bool {
	if plates[0] == result[0] && plates[1] == result[1] && plates[2] == result[2] {
		return true
	}

	return false
}

func isRespin(plates []int) bool {
	if plates[0] == 0 && plates[2] == 0 {
		return true
	}
	return false
}

func isSpWild(plates []int) bool {
	if plates[1] == wild2 || plates[1] == wild3 || plates[1] == wild4 {
		return true
	}
	return false
}
