package gamerule

import (
	"fmt"
	"testing"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gameplate"
)

func TestTime(test *testing.T) {
	fmt.Println(int(foundation.ServerNow().Weekday()))
}

func TestGetInitScroll(test *testing.T) {
	fmt.Println(foundation.JSONToString(GetInitBetRate()))
}

func Test243GameRequest(test *testing.T) {
	var att []map[string]interface{}
	var newatt []map[string]interface{}
	var result map[string]interface{}
	var otherdata map[string]interface{}
	var normalScore, scotterScore, totalwinscore, totalbetscore, betindex int64
	H5WildWinRate := make(map[int64]int64)
	normalPlateCount := make([]map[string]int64, len(scrollSize))
	normalSpWinCount := make(map[int64]int64)
	normalPlateScore := make(map[string]int)
	ScotterGameCount := make(map[int]int64)
	ScotterCombinationCount := make(map[string]int64)
	scotterPlateCount := make([]map[string]int64, len(scrollSize))
	scotterSpWinCount := make(map[int64]int64)
	scotterPlateScore := make(map[string]int)
	for i := range normalPlateCount {
		normalPlateCount[i] = make(map[string]int64)
	}
	for i := range scotterPlateCount {
		scotterPlateCount[i] = make(map[string]int64)
	}

	var playerid int64 = 24
	betmoney := GetBetMoney(betindex)
	attinfo := AttachInfo{PlayerID: playerid, Kind: GameIndex}
	att = attachInfoToAttachData(attinfo)

	for index, max := 0, 100000000; index < max; index++ {
		totalbetscore += betmoney

		result, newatt, otherdata = GameRequest(playerid, betindex, att)
		att[0] = newatt[0]
		if normalresult, ok := result["normalresult"]; ok {
			plate := ((normalresult).(map[string]interface{})["plate"]).([][]int)
			for i, rowSymbolarray := range plate {
				logplate := fmt.Sprintf("%v", rowSymbolarray)
				normalPlateCount[i][logplate]++
			}

			gameresult, ok := ((normalresult).(map[string]interface{})["gameresult"]).([]interface{})
			if ok && len(gameresult) > 0 {
				for i := range gameresult {
					InfoLine243Result := gameresult[i].(gameplate.InfoLine243)

					normalSpWinCount[InfoLine243Result.SpecialWinRate]++
					normalPlateScore[fmt.Sprintf("%v", InfoLine243Result.WinSymbolNum)] = InfoLine243Result.LineWinRate
				}
			}
		}
		normalTotalWin := result["totalwinscore"].(int64)
		normalScore += normalTotalWin
		totalwinscore += normalTotalWin

		scotteridarray := result["scotterid"].([]int64)
		aRoundScotterCount := 0
		for i, imax := 0, len(scotteridarray); i < imax; i++ {
			aRoundScotterCount++
			result, newatt, otherdata = ScotterGameRequest(playerid, betmoney, 6, newatt)
			scotteridarray := result["scotterid"].([]int64)
			ScotterCombinationCount[fmt.Sprintf("%v", result["scottercombination"])]++
			imax += len(scotteridarray)

			if scotterresult, ok := result["scotterresult"]; ok {
				for _, value := range scotterresult.([]interface{}) {

					H5Score := ((value).(map[string]interface{})["h5score"]).(int64)
					plate := ((value).(map[string]interface{})["plate"]).([][]int)
					for i, rowSymbolarray := range plate {
						logplate := fmt.Sprintf("%v", rowSymbolarray)
						scotterPlateCount[i][logplate]++
					}
					H5WildWinRate[H5Score/betmoney]++

					gameresults, ok := ((value).(map[string]interface{})["gameresult"]).([]interface{})
					if ok && len(result) > 0 {
						for _, result := range gameresults {
							InfoLine243Result := result.(gameplate.InfoLine243)
							SpWinRate := InfoLine243Result.SpecialWinRate

							scotterSpWinCount[SpWinRate]++
							scotterPlateScore[fmt.Sprintf("%v", InfoLine243Result.WinSymbolNum)] = InfoLine243Result.LineWinRate
						}
					}
				}
			}
			scotterWinScore := result["totalwinscore"].(int64)
			scotterScore += scotterWinScore
			totalwinscore += scotterWinScore
		}
		ScotterGameCount[aRoundScotterCount]++

		// moneypool.RTPControl(betmoney, otherdata["totalwinscore"])
		if otherdata["totalwinscore"].(int64) > 0 {
			// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		}
		if index%(max/100) == 0 {
			fmt.Println("Progress Rate:", float64(index)/float64(max)*100)
		}
	}

	fmt.Println("normalPlateCount:", foundation.JSONToString(normalPayCount))
	fmt.Println("normalPlateCount:", foundation.JSONToString(normalPlateCount))
	fmt.Println("normalSpWinCount:", foundation.JSONToString(normalSpWinCount))
	fmt.Println("normalPlateScore:", foundation.JSONToString(normalPlateScore))

	fmt.Println("normalPlateCount:", foundation.JSONToString(scotterPayCount))
	fmt.Println("scotterPlateCount:", foundation.JSONToString(scotterPlateCount))
	fmt.Println("scotterSpWinCount:", foundation.JSONToString(scotterSpWinCount))
	fmt.Println("scotterPlateScore:", foundation.JSONToString(scotterPlateScore))
	fmt.Println("ScotterGameCount:", foundation.JSONToString(ScotterGameCount))
	fmt.Println("ScotterCombinationCount:", foundation.JSONToString(ScotterCombinationCount))
	fmt.Println("H5WildWinRate:", foundation.JSONToString(H5WildWinRate))
	fmt.Println("Normal RTP:", fmt.Sprintf("%.2f", RTP(normalScore, totalbetscore)), "normalScore:", normalScore, "TotalBet:", totalbetscore)
	fmt.Println("Scotter RTP:", fmt.Sprintf("%.2f", RTP(scotterScore, totalbetscore)), "scotterScore:", scotterScore, "TotalBet:", totalbetscore)
	fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	// fmt.Println("InRespintCount:", inRespintCount, "RespinCountMap:", foundation.JSONToString(respinCountMap))

}

func Test243ScotterGameRequest(test *testing.T) {

}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
