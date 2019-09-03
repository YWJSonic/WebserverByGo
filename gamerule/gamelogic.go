package gamerule

import (
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/foundation/math"
	"gitlab.com/ServerUtility/gameplate"
)

// Result ...
func logicResult(betMoney int64, attinfo *AttachInfo) (map[string]interface{}, map[string]interface{}) {
	var result = make(map[string]interface{})
	var totalWin int64

	option := gameplate.PlateOption{
		Scotter: []int{scotter1, scotter2},
		Wild:    []int{wild1},
	}

	normalresult, otherdata, normaltotalwin := outputGame(betMoney, attinfo, option)
	FreeGameCount := foundation.InterfaceToInt(otherdata["freegamecount"])
	result["freegamecount"] = FreeGameCount
	result["normalresult"] = normalresult
	result["isfreegame"] = 0
	totalWin += normaltotalwin

	if FreeGameCount > 0 {
		freeresult, _, freetotalwin := outputFreeGame(betMoney, FreeGameCount, attinfo, option)
		result["freeresult"] = freeresult
		result["isfreegame"] = 1
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

	normalResult, otherdata, totalScores = aRound(betMoney, getNormalScorll(), option, 1)

	return normalResult, otherdata, totalScores
}

func outputFreeGame(betMoney int64, freeCount int, attinfo *AttachInfo, option gameplate.PlateOption) ([]map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	// var wildCount, bonusRate int
	otherdata := make(map[string]interface{})
	var freeResult []map[string]interface{}

	for i, imax := 0, freeCount; i < imax; i++ {
		tmpResult, _, tmpTotalScores := aRound(betMoney, getFreeScorll(), option, 2)
		totalScores += tmpTotalScores
		freeResult = append(freeResult, tmpResult)
	}
	return freeResult, otherdata, totalScores
}

func aRound(betMoney int64, scorll [][]int, option gameplate.PlateOption, gameType int) (map[string]interface{}, map[string]interface{}, int64) {

	var winLineInfo = make([]interface{}, 0)
	var totalScores int64
	var freeGameCount, scotterCount int
	// var paylinestr string
	var isLink bool
	var scotterLineSymbol, scotterLinePoint [][]int
	var plateSymbolCollectResult map[string]interface{}
	result := map[string]interface{}{
		"bonusrate":  int64(0),
		"bonusscore": int64(0),
	}
	otherdata := map[string]interface{}{
		"isfreegame":    0,
		"freegamecount": 0,
	}

	plateIndex, plateSymbol := gameplate.NewPlate2D(scrollSize, scorll)
	// plateSymbol = [][]int{
	// 	{5, 7, 8},
	// 	{2, 1, 2},
	// 	{1, 2, 1},
	// 	{1, 2, 1},
	// 	{7, 10, 8},
	// }
	plateLineMap := gameplate.PlateToLinePlate(plateSymbol, lineMap)

	for lineIndex, plateLine := range plateLineMap {
		newLine := gameplate.CutSymbolLink(plateLine, option) // cut line to win line point
		mulityLine := gameplate.LineMulitResult(newLine, option)

		if len(mulityLine) > 1 {
			isLink = false
			infoLine := gameplate.NewInfoLine()
			for _, winLine := range mulityLine {
				for _, payLine := range itemResults[len(winLine)] {
					if isWin(winLine, payLine, option) {
						isLink = true
						tmpline := winResult(betMoney, lineIndex, newLine, payLine, option, gameType)
						if tmpline.Score > infoLine.Score {
							infoLine = tmpline
							// paylinestr = fmt.Sprintf("%v", payLine[:len(payLine)-1])
						}
					}
				}
			}
			if isLink {
				totalScores += infoLine.Score
				winLineInfo = append(winLineInfo, infoLine)
			}
		} else {
			for _, payLine := range itemResults[len(newLine)] { // win line result group
				if isWin(newLine, payLine, option) { // win result check
					infoLine := winResult(betMoney, lineIndex, newLine, payLine, option, gameType)
					totalScores += infoLine.Score
					winLineInfo = append(winLineInfo, infoLine)
				}
			}
		}
	}

	// scotter 1 handle
	plateSymbolCollectResult = gameplate.PlateSymbolCollect(scotter1, plateSymbol, option, map[string]interface{}{"isincludewild": false, "isseachallplate": true})
	scotterCount = foundation.InterfaceToInt(plateSymbolCollectResult["targetsymbolcount"])
	scotterCount = math.ClampInt(scotterCount, 0, len(freeGameCountAttay))
	scotterLineSymbol = plateSymbolCollectResult["symbolnumcollation"].([][]int)
	scotterLinePoint = plateSymbolCollectResult["symbolpointcollation"].([][]int)

	if scotterCount >= scotter1GameLimit {
		infoLine := gameplate.NewInfoLine()

		for i, max := 0, len(scotterLineSymbol); i < max; i++ {
			if len(scotterLineSymbol[i]) > 0 {
				infoLine.AddNewLine(scotterLineSymbol[i], scotterLinePoint[i], option)
			} else {
				infoLine.AddEmptyPoint()
			}
		}

		infoLine.LineWinRate = scotter1LineRate[scotterCount]
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
		totalScores += infoLine.Score
		winLineInfo = append(winLineInfo, infoLine)

		freeGameCount = freeGameCountAttay[scotterCount]
		otherdata["freegamecount"] = freeGameCount
		otherdata["isfreegame"] = 1

	}

	// scotter 2 handle
	plateSymbolCollectResult = gameplate.PlateSymbolCollect(scotter2, plateSymbol, option, map[string]interface{}{"isincludewild": false, "isseachallplate": true})
	scotterCount = foundation.InterfaceToInt(plateSymbolCollectResult["targetsymbolcount"])
	scotterCount = math.ClampInt(scotterCount, 0, len(freeGameCountAttay))
	scotterLineSymbol = plateSymbolCollectResult["symbolnumcollation"].([][]int)
	scotterLinePoint = plateSymbolCollectResult["symbolpointcollation"].([][]int)

	if scotterCount >= scotter2GameLimit {
		infoLine := gameplate.NewInfoLine()

		for i, max := 0, len(scotterLineSymbol); i < max; i++ {
			if len(scotterLineSymbol[i]) > 0 {
				infoLine.AddNewLine(scotterLineSymbol[i], scotterLinePoint[i], option)
			} else {
				infoLine.AddEmptyPoint()
			}
		}

		infoLine.LineWinRate = scotter1LineRate[scotterCount]
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
		totalScores += infoLine.Score
		winLineInfo = append(winLineInfo, infoLine)

		bonusrate := bonusRate[foundation.RangeRandom(getScotter2Weightings())]
		result["bonusrate"] = bonusrate
		result["bonusscore"] = bonusrate * betMoney
		totalScores += bonusrate * betMoney

	}

	if len(winLineInfo) > 0 {
		result = foundation.AppendMap(result, gameplate.ResultMapLine(plateIndex, plateSymbol, winLineInfo))
	} else {
		result = foundation.AppendMap(result, gameplate.ResultMapLine(plateIndex, plateSymbol, []interface{}{}))
	}

	result["gameresult"] = winLineInfo
	result["scores"] = totalScores
	return result, otherdata, totalScores
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

func winResult(betMoney int64, lineIndex int, newLine, payLine []int, option gameplate.PlateOption, gameType int) gameplate.InfoLine {
	mainSymbol := payLine[0]
	infoLine := gameplate.NewInfoLine()

	for i, max := 0, len(payLine)-1; i < max; i++ {
		infoLine.AddNewPoint(newLine[i], lineMap[lineIndex][i], option)
	}

	if isScotter, _ := option.IsScotter(mainSymbol); isScotter {
		infoLine.LineWinRate = payLine[len(payLine)-1]
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
	} else {
		infoLine.LineWinRate = payLine[len(payLine)-1]
		infoLine.Score = int64(infoLine.LineWinRate) * (betMoney / betLine)
	}
	return infoLine
}
