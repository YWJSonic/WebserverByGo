package gamerule

import (
	"fmt"
	"testing"

	"gitlab.com/ServerUtility/gameplate"

	"gitlab.com/ServerUtility/foundation"
)

func TestTime(test *testing.T) {
	fmt.Println(int(foundation.ServerNow().Weekday()))
}

func TestGetInitScroll(test *testing.T) {
	fmt.Println(foundation.JSONToString(GetInitBetRate()))
}
func Test243GameRequest(test *testing.T) {
	var att []map[string]interface{}
	// var newatt []map[string]interface{}
	var result map[string]interface{}
	var otherdata map[string]interface{}
	var totalwinscore, totalbetscore, betindex int64
	var tmpLine [][]int
	var roundWinScore int64
	var ScotterGameCount int64
	normalSpWinCount := make(map[int64]int64)
	normalPlateCount := make([]map[string]int64, len(scrollSize))
	for i := range normalPlateCount {
		normalPlateCount[i] = make(map[string]int64)
	}
	// respinPlateCount := make(map[string]int64)
	normalPlateScore := make(map[string]int64)
	// respinPlateScore := make(map[string]int64)

	var playerid int64 = 24
	betmoney := GetBetMoney(betindex)
	attinfo := AttachInfo{PlayerID: playerid, Kind: GameIndex}
	att = attachInfoToAttachData(attinfo)

	for index, max := 0, 100000000; index < max; index++ {
		totalbetscore += betmoney
		roundWinScore = 0

		result, _, otherdata = GameRequest(playerid, 0, att)
		// fmt.Println(foundation.JSONToString(result))
		// att = newatt
		if normalresult, ok := result["normalresult"]; ok {
			gameresult := ((normalresult).(map[string]interface{})["gameresult"]).([]interface{})
			plate := ((normalresult).(map[string]interface{})["plate"]).([][]int)
			scores := ((normalresult).(map[string]interface{})["scores"]).(int64)
			for i, rowSymbolarray := range plate {
				logplate := fmt.Sprintf("%v", rowSymbolarray)
				normalPlateCount[i][logplate]++

			}

			if len(gameresult) > 0 {
				for _, gameWinLineInfo := range gameresult {
					WinLineInfo := (gameWinLineInfo).(gameplate.InfoLine243)
					roundWinScore += WinLineInfo.Score
					lineCount := 1
					for _, winSymbol := range WinLineInfo.LineSymbolNum {
						lineCount *= len(winSymbol)
					}

					tmpLine = make([][]int, lineCount)
					for i, imax := 0, len(WinLineInfo.LineSymbolNum); i < imax; i++ {
						for j, jmax := 0, lineCount; j < jmax; j++ {
							rowLimit := len(WinLineInfo.LineSymbolNum[i])
							tmpLine[j] = append(tmpLine[j], WinLineInfo.LineSymbolNum[i][j%rowLimit])
						}
					}
					for _, winline := range tmpLine {
						normalPlateScore[fmt.Sprintf("%v", winline)] = int64(WinLineInfo.LineWinRate / lineCount)
					}
				}
				if scores != roundWinScore {
					normalSpWinCount[(scores-roundWinScore)/betmoney]++
				}
			}
		}
		if otherdata["isscotter"].(int) == 1 {
			ScotterGameCount++
		}

		// moneypool.RTPControl(betmoney, otherdata["totalwinscore"])
		totalwinscore += otherdata["totalwinscore"].(int64)
		if otherdata["totalwinscore"].(int64) > 0 {
			// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		}
	}

	fmt.Println("ScotterGameCount", ScotterGameCount)
	fmt.Println("NormalPlateCount:", foundation.JSONToString(normalPlateCount))
	fmt.Println("normalSpWinCount:", foundation.JSONToString(normalSpWinCount))
	fmt.Println("normalPlateScore:", foundation.JSONToString(normalPlateScore))
	fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	// fmt.Println("InRespintCount:", inRespintCount, "RespinCountMap:", foundation.JSONToString(respinCountMap))

}

func Test243ScotterGameRequest(test *testing.T) {

}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
