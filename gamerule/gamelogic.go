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

	option := gameplate.PlateOption{
		Scotter: []int{scotter},
		Wild:    []int{wild1},
	}

	normalresult, otherdata, normaltotalwin := outputGame(betMoney, attinfo, option)
	result["normalresult"] = normalresult
	totalWin += normaltotalwin

	if iscotter, ok := otherdata["isscotter"]; ok && iscotter.(int) == 1 {
		freeresult, freeotherdata, freetotalwin := outputFreeGame(betMoney, attinfo, option)
		result["freeresult"] = freeresult
		result["freewildbonusrate"] = freeotherdata["wildbonusrate"]
		totalWin += freetotalwin
	}

	result["totalwinscore"] = totalWin
	return result, otherdata
}

// outputGame out put normal game result, mini game status, totalwin
func outputGame(betMoney int64, attinfo *AttachInfo, option gameplate.PlateOption) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	normalResult := make(map[string]interface{})
	otherdata := make(map[string]interface{})

	randWild := randWild()
	normalResult, otherdata, totalScores = aRound(betMoney, getNormalScorll(), randWild, option, 1)
	normalResult["randwild"] = randWild
	// normalResult["randwild"] = [][]int{}

	return normalResult, otherdata, totalScores
}

func outputFreeGame(betMoney int64, attinfo *AttachInfo, option gameplate.PlateOption) ([]map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var wildCount, bonusRate int
	otherdata := make(map[string]interface{})
	var freeResult []map[string]interface{}
	var lockWildarray = make([][]int, len(scrollSize))

	for i, imax := 0, freeCount; i < imax; i++ {
		tmpResult, _, tmpTotalScores := aRound(betMoney, freeScroll, lockWildarray, option, 2)
		totalScores += tmpTotalScores
		freeResult = append(freeResult, tmpResult)

		lockWildarray = lockWild(tmpResult["plate"].([][]int), lockWildarray, option)
	}
	for _, colArray := range lockWildarray {
		wildCount += len(colArray)
	}
	// freeWildCount[fmt.Sprintf("%v", wildCount)]++

	for limitIndex, limitCount := range wildBonusLimit {
		if wildCount >= limitCount {
			bonusRate = wildBonusRate[limitIndex]
		}
	}
	if bonusRate > 0 {
		freeWildBonusRateCount[fmt.Sprintf("%v", bonusRate)]++
		totalScores *= int64(bonusRate)
		otherdata["wildbonusrate"] = bonusRate
	} else {
		otherdata["wildbonusrate"] = 0
	}
	return freeResult, otherdata, totalScores
}

var normalPayLineCount map[string]int
var freePayLineCount map[string]int
var freeWildCount map[string]int
var freeWildBonusRateCount map[string]int

func aRound(betMoney int64, scorll [][]int, randWild [][]int, option gameplate.PlateOption, gameType int) (map[string]interface{}, map[string]interface{}, int64) {

	var winLineInfo []interface{}
	var totalScores int64
	otherdata := make(map[string]interface{})
	result := make(map[string]interface{})

	plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, scorll)

	// set random wild
	plateSymbol = setRandomWild(plateSymbol, randWild)
	plateLineMap := gameplate.PlateToLinePlate(plateSymbol, lineMap)

	for lineIndex, plateLine := range plateLineMap {
		newLine := gameplate.CutSymbolLink(plateLine, option) // cut line to win line point
		for _, payLine := range itemResults[len(newLine)] {   // win line result group
			if isWin(newLine, payLine, option) { // win result check
				// if gameType == 1 {
				// 	normalPayLineCount[fmt.Sprintf("%v", payLine)]++
				// } else {
				// 	freePayLineCount[fmt.Sprintf("%v", payLine)]++
				// }

				infoLine := gameplate.NewInfoLine()

				for i, max := 0, len(payLine)-1; i < max; i++ {
					infoLine.AddNewPoint(newLine[i], lineMap[lineIndex][i], option)
				}
				infoLine.LineWinRate = payLine[len(payLine)-1]
				infoLine.Score = int64(infoLine.LineWinRate) * (betMoney / betLine)
				totalScores += infoLine.Score
				winLineInfo = append(winLineInfo, infoLine)
			}
		}
	}

	result["scores"] = totalScores
	result["gameresult"] = winLineInfo
	if len(winLineInfo) > 0 {
		result = gameplate.ResultMapLine(plateIndex, plateSymbol, winLineInfo)
	} else {
		result = gameplate.ResultMapLine(plateIndex, plateSymbol, []interface{}{})
	}

	if isFreeGame(plateSymbol, option) {
		otherdata["isscotter"] = 1
	} else {
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

func setRandomWild(plateSymbol [][]int, randomWildPoint [][]int) [][]int {

	if len(randomWildPoint) <= 0 {
		return plateSymbol
	}

	for i, imax := 0, len(plateSymbol); i < imax; i++ {
		for j, jmax := 0, len(randomWildPoint[i]); j < jmax; j++ {
			plateSymbol[i][randomWildPoint[i][j]] = wild1
		}
	}
	return plateSymbol
}

func randWild() [][]int {
	var randwild = make([][]int, len(scrollSize))
	var randpoint []int
	var pointArray = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}

	wildCount := randWildCount[foundation.RangeRandom(randWildWeightings)]

	if wildCount <= 0 {
		return [][]int{}
	}

	randpoint = foundation.RandomMutily(pointArray, wildCount)
	// sort.Sort(randpoint)

	var col = 0
	for _, value := range randpoint {
		col = value / 3
		randwild[col] = append(randwild[col], value%3)
	}

	return randwild
}

// plateToLinePlate ...
func plateToLinePlate(plate [][]int, lineMap [][]int) [][]int {
	var plateLineMap [][]int
	var plateline []int

	for _, linePoint := range lineMap {
		plateline = []int{}
		for lineIndex, point := range linePoint {
			plateline = append(plateline, plate[lineIndex][point])
		}
		plateLineMap = append(plateLineMap, plateline)
	}

	return plateLineMap
}

// CutSymbolLink get line link array
func cutSymbolLink(symbolLine []int, option gameplate.PlateOption) []int {
	var newSymbolLine []int
	mainSymbol := symbolLine[0]

	for _, symbol := range symbolLine {
		if isWild, _ := option.IsWild(symbol); isWild {

		} else if isWild, _ := option.IsWild(mainSymbol); isWild {
			mainSymbol = symbol
		} else if symbol != mainSymbol {
			break
		}

		newSymbolLine = append(newSymbolLine, symbol)
	}

	return newSymbolLine
}

// isWin symbol line compar parline is win
func isWin(lineSymbol []int, payLineSymbol []int, option gameplate.PlateOption) bool {
	targetSymbol := 0
	isWin := true
	EmptyNum := option.EmptyNum()
	mainSymbol := EmptyNum

	for lineIndex, max := 0, len(payLineSymbol)-1; lineIndex < max; lineIndex++ {
		targetSymbol = lineSymbol[lineIndex]

		if isWild, _ := option.IsWild(targetSymbol); isWild {
			if mainSymbol == EmptyNum {
				mainSymbol = targetSymbol
			}
			continue
		}

		switch payLineSymbol[lineIndex] {
		case targetSymbol:
			mainSymbol = targetSymbol
		default:
			isWin = false
			return isWin
		}
	}

	if mainSymbol != payLineSymbol[0] {
		return false
	}

	return isWin
}

func isFreeGame(plate [][]int, option gameplate.PlateOption) bool {
	var scotterCount = 0

	for _, colarray := range plate {
		for _, rowSymbol := range colarray {
			if isScotter, _ := option.IsScotter(rowSymbol); isScotter {
				scotterCount++
			}
		}
	}

	if scotterCount >= scotterGameLimit {
		return true
	}
	return false
}

func lockWild(plater [][]int, lockWild [][]int, option gameplate.PlateOption) [][]int {

	for colIndex, colarray := range plater {
		for rowIndex, row := range colarray {
			if isWild, _ := option.IsWild(row); isWild && !foundation.IsInclude(rowIndex, lockWild[colIndex]) {
				lockWild[colIndex] = append(lockWild[colIndex], rowIndex)
			}
		}
	}

	return lockWild
}
