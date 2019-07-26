package gamerule

import (
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gameplate"
)

// Result ...
func logicResult(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}) {
	var result = make(map[string]interface{})
	var totalWin int64

	normalresult, otherdata, normaltotalwin := outputGame(betMoney, attinfo)
	result = foundation.AppendMap(result, otherdata)
	result["normalresult"] = normalresult
	totalWin += normaltotalwin

	result["totalwinscore"] = totalWin
	return result, otherdata
}
func logicScotterGameResult(betMoney int64, scotterWinRateIndex, scotterSpinTimeIndex int, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}) {
	var result = make(map[string]interface{})
	var totalWin int64

	scotterresult, otherdata, scottertotalwin := outputFreeGame(betMoney, scotterWinRateIndex, scotterSpinTimeIndex, attinfo)
	result = foundation.AppendMap(result, otherdata)
	result["scotterresult"] = scotterresult
	totalWin += scottertotalwin

	result["totalwinscore"] = totalWin
	return result, otherdata
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var lineInfo gameplate.InfoLine243
	var winLineInfo []interface{}
	IsSpecialWin := false
	result := make(map[string]interface{})
	otherdata := make(map[string]interface{})
	otherdata["isscotter"] = 0
	option := gameplate.PlateOption{
		Scotter: []int{scotter},
		Wild:    []int{wild1},
	}
	plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, normalScroll)

	for _, ItemNum := range items {
		mainSymbol, symbolNumCollation, symBolPointCollation := symbolCollation(ItemNum, plateSymbol, option)

		for _, payLine := range itemResults[len(symbolNumCollation)] {
			if mainSymbol == payLine[0] {
				lineInfo = newBaseInfoLine(symbolNumCollation, symBolPointCollation, payLine, betMoney, option)
				totalScores += lineInfo.Score
				winLineInfo = append(winLineInfo, lineInfo)

				if lineInfo.WildCount() > 0 {
					IsSpecialWin = true
				}
			}
		}
	}
	if scotterCount(plateSymbol, option) >= scotterGameLimit {
		otherdata["isscotter"] = 1
	}
	if IsSpecialWin {
		var H5RandWinRate int
		H5RandWinRate = normalWildWinRate[foundation.RangeRandom(normalWildWinRateWeightings)]
		totalScores += int64(H5RandWinRate) * betMoney
	}
	result = gameplate.ResultMap243(plateIndex, plateSymbol, winLineInfo)
	return result, otherdata, totalScores
}

func outputFreeGame(betMoney int64, scotterWinRateIndex, scotterSpinTimeIndex int, attinfo *AttachInfo) (interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var weekDay int64
	var lineInfo gameplate.InfoLine243
	var winLineInfo []interface{}
	var scotterResult []interface{}
	otherdata := make(map[string]interface{})
	otherdata["isscotter"] = 0
	otherdata["scottercount"] = 0
	IsSpecialWin := false
	WinRate := scotterGameWildWinRate[scotterWinRateIndex]
	WinRateWeightings := scotterGameWildWinRateWeightings[scotterWinRateIndex]
	SpinTime := scotterGameSpinTime[scotterWinRateIndex]
	option := gameplate.PlateOption{
		Scotter: []int{scotter},
		Wild:    []int{wild1},
	}

	for i, imax := 0, SpinTime; i < imax; i++ {
		IsSpecialWin = false

		plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, normalScroll)
		for _, ItemNum := range items {
			mainSymbol, symbolNumCollation, symBolPointCollation := symbolCollation(ItemNum, plateSymbol, option)

			for _, payLine := range itemResults[len(symbolNumCollation)] {
				if mainSymbol == payLine[0] {
					lineInfo = newBaseInfoLine(symbolNumCollation, symBolPointCollation, payLine, betMoney, option)
					totalScores += lineInfo.Score
					winLineInfo = append(winLineInfo, lineInfo)

					if lineInfo.WildCount() > 0 {
						IsSpecialWin = true
					}
				}
			}
		}
		if scotterCount(plateSymbol, option) >= scotterGameLimit {

			weekDay = int64(foundation.ServerNow().Weekday())
			attinfo.DayScotterGameCount++
			ref := attinfo.DayScotterGameCount*100 + weekDay
			attinfo.DayScotterGameInfo[ref+10] = 0
			attinfo.FreeGameBetLockIndex[ref+20] = betMoney
		}
		if IsSpecialWin {
			totalScores += WinRate[foundation.RangeRandom(WinRateWeightings)] * betMoney
		}

		otherdata["scottercount"] = 0
		scotterResult = append(scotterResult, gameplate.ResultMap243(plateIndex, plateSymbol, winLineInfo))
	}
	return scotterResult, otherdata, totalScores
}

func aRound(betMoney int64) (map[string]interface{}, map[string]interface{}, int64) {

	var winLineInfo []interface{}
	var IsSpecialWin bool
	var totalScores int64
	var lineInfo gameplate.InfoLine243
	otherdata := make(map[string]interface{})
	result := make(map[string]interface{})

	option := gameplate.PlateOption{
		Scotter: []int{scotter},
		Wild:    []int{wild1},
	}
	plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, normalScroll)

	for _, ItemNum := range items {
		mainSymbol, symbolNumCollation, symBolPointCollation := symbolCollation(ItemNum, plateSymbol, option)

		for _, payLine := range itemResults[len(symbolNumCollation)] {
			if mainSymbol == payLine[0] {
				lineInfo = newBaseInfoLine(symbolNumCollation, symBolPointCollation, payLine, betMoney, option)
				totalScores += lineInfo.Score
				winLineInfo = append(winLineInfo, lineInfo)

				if lineInfo.WildCount() > 0 {
					IsSpecialWin = true
				}
			}
		}
	}
	if scotterCount(plateSymbol, option) >= scotterGameLimit {
		otherdata["isscotter"] = 1
	}
	if IsSpecialWin {
		otherdata["isspecialwin"] = 1
	}
	result = gameplate.ResultMap243(plateIndex, plateSymbol, winLineInfo)
	return result, otherdata, totalScores
}

func symbolCollation(symbolNum int, plate [][]int, option gameplate.PlateOption) (int, [][]int, [][]int) {
	var symBolPointCollation [][]int
	var symbolNumCollation [][]int
	var mainSymbol = option.EmptyNum()

	for _, colArray := range plate {
		var rowPointArray []int
		var rowSymbolArray []int
		for rowIndex, rowSymbol := range colArray {
			if IsWild, _ := option.IsWild(rowSymbol); IsWild {
				rowSymbolArray = append(rowSymbolArray, rowSymbol)
				rowPointArray = append(rowPointArray, rowIndex)

				if isWild, _ := option.IsWild(symbolNum); isWild {
					mainSymbol = rowSymbol
				}

			} else if IsScotter, _ := option.IsScotter(rowSymbol); IsScotter {

			} else if symbolNum == rowSymbol {
				rowSymbolArray = append(rowSymbolArray, rowSymbol)
				rowPointArray = append(rowPointArray, rowIndex)
				mainSymbol = rowSymbol
			}
		}

		if len(rowPointArray) <= 0 {
			break
		}
		symbolNumCollation = append(symbolNumCollation, rowSymbolArray)
		symBolPointCollation = append(symBolPointCollation, rowPointArray)
	}

	if mainSymbol != symbolNum {
		return mainSymbol, make([][]int, 0), make([][]int, 0)
	}
	return mainSymbol, symbolNumCollation, symBolPointCollation
}

func newBaseInfoLine(lineSymbol [][]int, linePoint [][]int, payLine []int, betMoney int64, option gameplate.PlateOption) gameplate.InfoLine243 {
	infoLine := gameplate.NewInfoLine243()

	for i := range lineSymbol {
		infoLine.AddNewLine(lineSymbol[i], linePoint[i], option)
	}

	infoLine.LineWinRate = payLine[len(payLine)-1]
	infoLine.Score = int64(infoLine.LineWinRate) * betMoney
	infoLine.IsLink = 1
	return infoLine
}

func mysteryCombination() []int {
	mysteryIndex := foundation.RangeRandom(scotterGameMysteryWeightings)
	mysteryIndexCombination := scotterGameMysteryIndexCombination[mysteryIndex]
	return mysteryIndexCombination
}

func scotterCount(plate [][]int, option gameplate.PlateOption) int {
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
