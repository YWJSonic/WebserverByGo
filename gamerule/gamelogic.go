package gamerule

import (
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
	}
	_, plate := gameplate.NewPlate2D(scrollSize, normalScroll)

	// var symBolCollation [][]int

	for _, ItemNum := range items {
		symbolNumCollation, symBolPointCollation := symbolCollation(ItemNum, plate, option)

		for _, payLine := range itemResults[len(symbolNumCollation)] {
			if isWin(symbolNumCollation, payLine, option) {
				for plateIndex := range symbolNumCollation {
					infoLineAddNewPoint(symbolNumCollation[plateIndex], symBolPointCollation[plateIndex], payLine, option)
				}
			}
		}

	}

	result["scores"] = totalScores
	result["gameresult"] = winLineInfo
	return result, otherdata, totalScores
}

func symbolCollation(symbolNum int, plate [][]int, option gameplate.PlateOption) ([][]int, [][]int) {
	var symBolPointCollation [][]int
	var symbolNumCollation [][]int
	var MainSymbol = option.EmptyNum()

	for _, colArray := range plate {
		var rowPointArray []int
		var rowSymbolArray []int
		for rowIndex, rowSymbol := range colArray {
			if IsWild, _ := option.IsWild(rowSymbol); IsWild {
				rowSymbolArray = append(rowSymbolArray, rowSymbol)
				rowPointArray = append(rowPointArray, rowIndex)

				if isWild, _ := option.IsWild(symbolNum); isWild {
					MainSymbol = rowSymbol
				}

			} else if IsScotter, _ := option.IsScotter(rowSymbol); IsScotter {

			} else if symbolNum == rowSymbol {
				rowSymbolArray = append(rowSymbolArray, rowSymbol)
				rowPointArray = append(rowPointArray, rowIndex)
				MainSymbol = rowSymbol
			}
		}

		if len(rowPointArray) <= 0 {
			break
		}
		symbolNumCollation = append(symbolNumCollation, rowSymbolArray)
		symBolPointCollation = append(symBolPointCollation, rowPointArray)
	}

	if MainSymbol != symbolNum {
		return make([][]int, 0), make([][]int, 0)
	}
	return symbolNumCollation, symBolPointCollation
}

// isWin symbol line compar parline is win
func isWin(lineSymbol [][]int, payLineSymbol []int, option gameplate.PlateOption) bool {
	// targetSymbol := 0
	isWin := true

	// for lineIndex, max := 0, len(payLineSymbol)-1; lineIndex < max; lineIndex++ {
	// 	targetSymbol = lineSymbol[lineIndex]

	// 	if isWild, _ := option.IsWild(targetSymbol); isWild {
	// 		continue
	// 	}

	// 	switch payLineSymbol[lineIndex] {
	// 	case targetSymbol:
	// 	case -1000:
	// 		if !foundation.IsInclude(targetSymbol, symbolGroup[-1000]) {
	// 			isWin = false
	// 			return isWin
	// 		}
	// 	case -1001:
	// 		if !foundation.IsInclude(targetSymbol, symbolGroup[-1001]) {
	// 			isWin = false
	// 			return isWin
	// 		}
	// 	default:
	// 		isWin = false
	// 		return isWin
	// 	}
	// }
	return isWin
}

func infoLineAddNewPoint(lineSymbol []int, linePoint []int, lineWinResult []int, option gameplate.PlateOption) gameplate.InfoLine {
	infoLine := gameplate.NewInfoLine()

	for i, max := 0, len(lineWinResult)-1; i < max; i++ {
		infoLine.AddNewPoint(lineSymbol[i], linePoint[i], option)
	}

	return infoLine
}

// func wildCount()int,[][]int{
// }

func processInfoLine(betMoney int64, winLineInfo *gameplate.InfoLine, option gameplate.PlateOption) {

	// if winLineInfo.WinRate > 0 {
	// 	winLineInfo.Score = int64(winLineInfo.WinRate) * betMoney
	// } else {
	// 	switch winLineInfo.WinRate {
	// 	case -100:
	// 		for _, payLine := range itemResults[len(winLineInfo.LineSymbolNum)] {
	// 			if isWin(payLine, []int{winLineInfo.LineSymbolNum[0][0], winLineInfo.LineSymbolNum[0][0], winLineInfo.LineSymbolNum[0][0]}, option) {
	// 				winLineInfo.WinRate = payLine[len(payLine)-1]
	// 			}
	// 		}
	// 	case -101:
	// 		winLineInfo.WinRate = 30
	// 	case -102:
	// 		winLineInfo.WinRate = 45
	// 	case -103:
	// 		winLineInfo.WinRate = 75
	// 	}
	// 	winLineInfo.Score = int64(winLineInfo.WinRate) * betMoney
	// }
}

// func winLineInfoAnalysis(winLineInfo []int) {
// }
