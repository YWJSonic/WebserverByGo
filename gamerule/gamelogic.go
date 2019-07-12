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

	// if otherdata["isrespin"].(int) == 1 {
	// 	respinresult, respintotalwin := outRespin(betMoney, attinfo)
	// 	totalWin += respintotalwin
	// 	result["respin"] = respinresult
	// 	result["isrespin"] = 1
	// }

	result["totalwinscore"] = totalWin
	return result
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var winLineInfo []interface{}
	result := make(map[string]interface{})
	otherdata := make(map[string]interface{})
	option := gameplate.PlateOption{
		Wild:          []int{wild1},
		LineMiniCount: 3,
		Group:         symbolGroup,
	}
	_, plate := gameplate.NewPlate2D(scrollSize, normalScroll)
	plateLineMap := gameplate.PlateToLinePlate(plate, lineMap)

	for lineIndex, plateLine := range plateLineMap {
		newLine := gameplate.CutSymbolLink(plateLine, option)
		for _, payLine := range itemResults[len(newLine)] {
			if isWin(newLine, payLine, option) {
				totalScores += betMoney * int64(payLine[len(payLine)-1])
				fmt.Println("Win", newLine, payLine, lineMap[lineIndex], totalScores)
			}
		}
	}
	// infoLineAddNewPoint(plateLine, lineMap[lineIndex], payLine, option)
	result["scores"] = totalScores
	result["gameresult"] = winLineInfo
	return result, otherdata, totalScores
}

func infoLineAddNewPoint(symbol []int, linePoint []int, lineWinResult []int, option gameplate.PlateOption) gameplate.InfoLine {
	infoLine := gameplate.NewLineInfo()

	for i, max := 0, len(lineWinResult)-1; i < max; i++ {
		infoLine.AddNewPoint(symbol[i], linePoint[i], option)
	}

	return infoLine
}

// isWin symbol line compar parline is win
func isWin(lineSymbol []int, payLineSymbol []int, option gameplate.PlateOption) bool {
	targetSymbol := 0
	isWin := true

	for lineIndex, max := 0, len(payLineSymbol)-1; lineIndex < max; lineIndex++ {
		targetSymbol = lineSymbol[lineIndex]

		if isWild, _ := option.IsWild(targetSymbol); isWild {
			continue
		}

		switch payLineSymbol[lineIndex] {
		case targetSymbol:
		case -1000:
			if !foundation.IsInclude(targetSymbol, symbolGroup[-1000]) {
				isWin = false
				return isWin
			}
		case -1001:
			if !foundation.IsInclude(targetSymbol, symbolGroup[-1001]) {
				isWin = false
				return isWin
			}
		default:
			isWin = false
			return isWin
		}
	}
	return isWin
}
