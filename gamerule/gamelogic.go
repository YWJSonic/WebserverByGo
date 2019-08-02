package gamerule

import (
	"fmt"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gameplate"
)

// Result ...
func logicResult(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}) {
	var result = make(map[string]interface{})
	var totalWin int64
	var scotterIDArray []int64

	normalresult, otherdata, normaltotalwin := outputGame(betMoney, attinfo)
	result["normalresult"] = normalresult
	totalWin += normaltotalwin

	if otherdata["isscotter"] == 1 {
		scotterid := attinfo.NewDayScotterID()
		attinfo.DayScotterGameInfo[DayScotterGameInfoKey(scotterid)] = 0
		attinfo.FreeGameBetLockIndex[FreeGameBetLockIndexKey(scotterid)] = betMoney
		scotterIDArray = append(scotterIDArray, scotterid)
		result["scotterid"] = scotterIDArray
	} else {
		result["scotterid"] = []int64{}
	}

	result["totalwinscore"] = totalWin
	return result, otherdata
}
func logicScotterGameResult(betMoney int64, scotterWinRateIndex, scotterSpinTimeIndex int, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}) {
	var result = make(map[string]interface{})
	var totalWin int64
	var scotterIDArray []int64

	scotterresult, otherdata, scottertotalwin := outputScotterGame(betMoney, scotterWinRateIndex, scotterSpinTimeIndex, attinfo)
	// result = foundation.AppendMap(result, otherdata)
	result["scotterresult"] = scotterresult
	totalWin += scottertotalwin
	scotterCount := otherdata["scottercount"].(int)

	if otherdata["isscotter"] == 1 {
		for i, imax := 0, scotterCount; i < imax; i++ {
			scotterid := attinfo.NewDayScotterID()
			attinfo.DayScotterGameInfo[DayScotterGameInfoKey(scotterid)] = 0
			attinfo.FreeGameBetLockIndex[FreeGameBetLockIndexKey(scotterid)] = betMoney
			scotterIDArray = append(scotterIDArray, scotterid)
		}
		result["scotterid"] = scotterIDArray
	} else {
		result["scotterid"] = []int64{}

	}

	result["totalwinscore"] = totalWin
	return result, otherdata
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	// var scotterIDArray []int64
	normalResult := make(map[string]interface{})
	otherdata := make(map[string]interface{})

	normalResult, otherdata, totalScores = aRound(betMoney, normalWildWinRate, normalWildWinRateWeightings, normalScroll)

	return normalResult, otherdata, totalScores
}

func outputScotterGame(betMoney int64, scotterWinRateIndex, scotterSpinTimeIndex int, attinfo *AttachInfo) (interface{}, map[string]interface{}, int64) {
	var scotterCount int
	var totalScores, tmpScores int64
	var scotterResult []interface{}
	var tmpResult map[string]interface{}
	tmpOtherdata := make(map[string]interface{})
	otherdata := make(map[string]interface{})
	WinRate := scotterGameWildWinRate[scotterWinRateIndex]
	WinRateWeightings := scotterGameWildWinRateWeightings[scotterWinRateIndex]
	SpinTime := scotterGameSpinTime[scotterWinRateIndex]

	for i, imax := 0, SpinTime; i < imax; i++ {

		tmpResult, tmpOtherdata, tmpScores = aRound(betMoney, WinRate, WinRateWeightings, scotterScroll)
		if tmpOtherdata["isscotter"] == 1 {
			scotterCount++
		}

		if isSpeicalH5Win(tmpResult["plate"].([][]int)) {
			totalScores += int64(scotterH5SpecialWinRate[foundation.RangeRandom(scotterH5SpecialWinRateWeightings[scotterWinRateIndex])]) * betMoney
			tmpResult["h5score"] = int64(scotterH5SpecialWinRate[foundation.RangeRandom(scotterH5SpecialWinRateWeightings[scotterWinRateIndex])]) * betMoney
		} else {
			tmpResult["h5score"] = int64(0)
		}
		totalScores += tmpScores
		scotterResult = append(scotterResult, tmpResult)
	}

	if scotterCount > 0 {
		otherdata["isscotter"] = 1
	} else {
		otherdata["isscotter"] = 0
	}
	otherdata["scottercount"] = scotterCount
	return scotterResult, otherdata, totalScores
}

var normalPayCount map[string]int64
var scotterPayCount map[string]int64

func init() {
	normalPayCount = make(map[string]int64)
	scotterPayCount = make(map[string]int64)
}

func aRound(betMoney int64, spWinRate, spWinWeightings []int64, scorll [][]int) (map[string]interface{}, map[string]interface{}, int64) {

	var winLineInfo []interface{}
	var IsScotter bool
	var wildRandWinRate int64
	var totalScores int64
	var lineInfo gameplate.InfoLine243
	otherdata := make(map[string]interface{})
	result := make(map[string]interface{})

	option := gameplate.PlateOption{
		Scotter: []int{scotter},
		Wild:    []int{wild1},
	}
	plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, scorll)
	plateSymbol = [][]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 5, 6},
		{1, 5, 6},
		{1, 5, 6},
	}
	for _, ItemNum := range items {
		symbolNumCollation, symBolPointCollation := symbolCollation(ItemNum, plateSymbol, option)

		IsScotterSymbol, _ := option.IsScotter(ItemNum)
		if IsScotterSymbol {
			if len(symbolNumCollation) >= scotterGameLimit {
				IsScotter = true
			}
		}

		for _, payLine := range itemResults[len(symbolNumCollation)] {
			if ItemNum == payLine[0] {
				wildRandWinRate = spWinRate[foundation.RangeRandomInt64(spWinWeightings)]
				// lineInfo = newBaseInfoLine(symbolNumCollation, symBolPointCollation, payLine, betMoney, wildRandWinRate, option)
				var lineCount, lineNormal = 1, 1
				lineInfo = gameplate.NewInfoLine243()

				for i := range symbolNumCollation {
					lineInfo.AddNewWinSymol(symbolNumCollation[i], symBolPointCollation[i], option)
					lineCount *= len(symbolNumCollation[i])
					lineNormal *= len(symbolNumCollation[i]) - len(lineInfo.WildPoint[i])
				}

				if IsScotterSymbol {
					lineInfo.LineWinRate = payLine[len(payLine)-1]
					lineInfo.Score = int64(lineInfo.LineWinRate) * betMoney * int64(lineCount)
				} else {
					if (lineCount - lineNormal) > 0 {
						lineInfo.LineWinRate = payLine[len(payLine)-1]
						lineInfo.SpecialWinRate = wildRandWinRate
						lineInfo.Score += int64(lineNormal) * int64(lineInfo.LineWinRate) * (betMoney / betLine)
						lineInfo.Score += int64(lineCount-lineNormal) * int64(lineInfo.LineWinRate) * (betMoney / betLine) * wildRandWinRate
					} else {
						lineInfo.LineWinRate = payLine[len(payLine)-1]
						lineInfo.Score = int64(lineInfo.LineWinRate) * (betMoney / betLine) * int64(lineCount)
					}
				}
				totalScores += lineInfo.Score
				winLineInfo = append(winLineInfo, lineInfo)

				if len(spWinRate) > 3 {
					normalPayCount[fmt.Sprintf("%v", payLine)] += int64(lineCount)
				} else {
					scotterPayCount[fmt.Sprintf("%v", payLine)] += int64(lineCount)
				}
			}
		}
	}

	if len(winLineInfo) > 0 {
		result = gameplate.ResultMap243(plateIndex, plateSymbol, winLineInfo)
	} else {
		result = gameplate.ResultMap243(plateIndex, plateSymbol, []interface{}{})
	}

	if IsScotter {
		result["isscotter"] = 1
		otherdata["isscotter"] = 1
	} else {
		result["isscotter"] = 0
		otherdata["isscotter"] = 0
	}
	return result, otherdata, totalScores
}

func symbolCollation(symbolNum int, plate [][]int, option gameplate.PlateOption) ([][]int, [][]int) {
	var symBolPointCollation = make([][]int, 0)
	var symbolNumCollation = make([][]int, 0)
	// var mainSymbol = option.EmptyNum()
	var IsWildTarget, _ = option.IsWild(symbolNum)
	var IsScotterTarget, _ = option.IsScotter(symbolNum)

	for _, colArray := range plate {
		var rowPointArray []int
		var rowSymbolArray []int
		for rowIndex, rowSymbol := range colArray {
			if IsWildTarget {
				if symbolNum == rowSymbol {
					rowSymbolArray = append(rowSymbolArray, rowSymbol)
					rowPointArray = append(rowPointArray, rowIndex)
					// mainSymbol = rowSymbol
				}
			} else if IsScotterTarget {
				if symbolNum == rowSymbol {
					rowSymbolArray = append(rowSymbolArray, rowSymbol)
					rowPointArray = append(rowPointArray, rowIndex)
					// mainSymbol = rowSymbol
				}
			} else {
				IsWild, _ := option.IsWild(rowSymbol)

				if symbolNum == rowSymbol {
					rowSymbolArray = append(rowSymbolArray, rowSymbol)
					rowPointArray = append(rowPointArray, rowIndex)
					// mainSymbol = rowSymbol
				} else if IsWild {
					rowSymbolArray = append(rowSymbolArray, rowSymbol)
					rowPointArray = append(rowPointArray, rowIndex)

					// if isWild, _ := option.IsWild(symbolNum); isWild {
					// 	mainSymbol = rowSymbol
					// }

				}
			}

		}

		if len(rowPointArray) <= 0 {
			break
		}
		symbolNumCollation = append(symbolNumCollation, rowSymbolArray)
		symBolPointCollation = append(symBolPointCollation, rowPointArray)
	}

	// if mainSymbol != symbolNum {
	// 	return make([][]int, 0), make([][]int, 0)
	// }
	return symbolNumCollation, symBolPointCollation
	// if mainSymbol != symbolNum {
	// 	return mainSymbol, make([][]int, 0), make([][]int, 0)
	// }
	// return mainSymbol, symbolNumCollation, symBolPointCollation
}

func newBaseInfoLine(lineSymbol [][]int, linePoint [][]int, payLine []int, betMoney, wildRandWinRate int64, option gameplate.PlateOption) gameplate.InfoLine243 {
	var lineCount, lineNormal = 1, 1
	lineInfo := gameplate.NewInfoLine243()

	for i := range lineSymbol {
		lineInfo.AddNewWinSymol(lineSymbol[i], linePoint[i], option)
		lineCount *= len(lineSymbol[i])
		lineNormal *= len(lineSymbol[i]) - len(lineInfo.WildPoint[i])
	}

	if (lineCount - lineNormal) > 0 {
		lineInfo.LineWinRate = payLine[len(payLine)-1]
		lineInfo.SpecialWinRate = wildRandWinRate
		lineInfo.Score += int64(lineNormal) * int64(lineInfo.LineWinRate) * (betMoney / betLine)
		lineInfo.Score += int64(lineCount-lineNormal) * int64(lineInfo.LineWinRate) * (betMoney / betLine) * wildRandWinRate
		fmt.Println(lineSymbol, payLine[len(payLine)-1], wildRandWinRate, lineInfo.Score)
	} else {
		lineInfo.LineWinRate = payLine[len(payLine)-1]
		lineInfo.Score = int64(lineInfo.LineWinRate) * (betMoney / betLine) * int64(lineCount)
		fmt.Println(lineSymbol, payLine[len(payLine)-1], 0, lineInfo.Score)
	}
	return lineInfo
}
func mysteryCombination() []int {
	mysteryIndex := foundation.RangeRandom(scotterGameMysteryWeightings)
	mysteryIndexCombination := scotterGameMysteryIndexCombination[mysteryIndex]
	return mysteryIndexCombination
}

func plateScotterCount(plate [][]int, option gameplate.PlateOption) int {
	var scotterCount int
	for _, col := range plate {
		for _, row := range col {
			if isScotter, _ := option.IsScotter(row); isScotter {
				scotterCount++
			}
		}
	}
	return scotterCount
}

func isSpeicalH5Win(plateSymbol [][]int) bool {
	var col1, col5 bool

	for _, rowSymbol := range plateSymbol[0] {
		if rowSymbol == scotterH5 {
			col1 = true
		}
	}
	for _, rowSymbol := range plateSymbol[4] {
		if rowSymbol == scotterH5 {
			col5 = true
		}
	}

	if col1 && col5 {
		return true
	}
	return false
}
