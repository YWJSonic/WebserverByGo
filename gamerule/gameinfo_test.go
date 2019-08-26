package gamerule

import (
	"fmt"
	"math/rand"
	"testing"

	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/ServerUtility/gameplate"
)

func TestTime(test *testing.T) {
	fmt.Println(int(foundation.ServerNow().Weekday()))
}

func TestGetInitScroll(test *testing.T) {

	for {
		println(rand.Intn(2))
	}
}

func Test243GameRequest(test *testing.T) {

	// for ii, iimax := 1, 6; ii < iimax; ii++ {
	var att []map[string]interface{}
	var result map[string]interface{}
	var otherdata map[string]interface{}
	var normalScore, freeScore, totalwinscore, totalbetscore, betindex int64

	randwildCount := make(map[string]int64)
	normalPayLineCount = make(map[string]int)
	freePayLineCount = make(map[string]int)
	freeWildCount = make(map[string]int)
	freeWildBonusRateCount = make(map[string]int)

	normalPlateCount := make([]map[string]int64, len(scrollSize))
	normalPlateScore := make(map[string]int)
	normalScoreCount := make(map[string]int)

	freePlateCount := make([]map[string]int64, len(scrollSize))
	freePlateScore := make(map[string]int)
	freeScoreCount := make(map[string]int)

	for i := range normalPlateCount {
		normalPlateCount[i] = make(map[string]int64)
		freePlateCount[i] = make(map[string]int64)
	}

	var playerid int64 = 24
	betmoney := GetBetMoney(betindex)
	attinfo := AttachInfo{PlayerID: playerid, Kind: GameIndex}
	att = attachInfoToAttachData(attinfo)
	RTPSetting = 5 // ii

	for index, max := 0, 1000000; index < max; index++ {
		totalbetscore += betmoney

		result, _, otherdata = GameRequest(playerid, betindex, att)
		if normalresult, ok := result["normalresult"]; ok {
			plate := ((normalresult).(map[string]interface{})["plate"]).([][]int)
			for i, rowSymbolarray := range plate {
				logplate := fmt.Sprintf("%v", rowSymbolarray)
				normalPlateCount[i][logplate]++
			}

			randwild := ((normalresult).(map[string]interface{})["randwild"]).([][]int)
			var wildcount int
			for _, col := range randwild {
				wildcount += len(col)
			}
			randwildCount[fmt.Sprintf("%v", wildcount)]++

			gameresult, ok := ((normalresult).(map[string]interface{})["gameresult"]).([]interface{})
			if ok && len(gameresult) > 0 {
				for _, value := range gameresult {
					roundResult := value.(gameplate.InfoLine)
					score := roundResult.Score
					normalScore += score
					totalwinscore += score
					normalPlateScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.LineWinRate
					normalScoreCount[fmt.Sprintf("%v", roundResult.LineWinRate)]++
				}
			}
		}
		// free game
		if freeresult, ok := result["freeresult"]; ok {
			result := freeresult.([]map[string]interface{})

			for _, value := range result {
				plate := (value["plate"]).([][]int)
				for i, rowSymbolarray := range plate {
					logplate := fmt.Sprintf("%v", rowSymbolarray)
					freePlateCount[i][logplate]++
				}

				gameresult, ok := (value["gameresult"]).([]interface{})
				if ok && len(gameresult) > 0 {
					for _, value := range gameresult {
						roundResult := value.(gameplate.InfoLine)
						score := roundResult.Score
						freeScore += score
						totalwinscore += score
						freePlateScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.LineWinRate
						freeScoreCount[fmt.Sprintf("%v", roundResult.LineWinRate)]++
					}
				}
			}
		}

		// moneypool.RTPControl(betmoney, otherdata["totalwinscore"])
		if otherdata["totalwinscore"].(int64) > 0 {
			// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		}
		if index%(max/100) == 0 {
			fmt.Println("Progress Rate:", float64(index)/float64(max)*100)
		}
	}

	fmt.Println("randwildCount:", foundation.JSONToString(randwildCount))

	fmt.Println("normalPlateCount:", foundation.JSONToString(normalPlateCount))
	fmt.Println("normalPlateScore:", foundation.JSONToString(normalPlateScore))
	fmt.Println("normalScoreCount:", foundation.JSONToString(normalScoreCount))
	fmt.Println("normalPayLineCount:", foundation.JSONToString(normalPayLineCount))

	fmt.Println("freePlateCount:", foundation.JSONToString(freePlateCount))
	fmt.Println("freePlateScore:", foundation.JSONToString(freePlateScore))
	fmt.Println("freeScoreCount:", foundation.JSONToString(freeScoreCount))
	fmt.Println("freePayLineCount:", foundation.JSONToString(freePayLineCount))
	fmt.Println("freeWildCount:", foundation.JSONToString(freeWildCount))
	fmt.Println("freeWildBonusRateCount:", foundation.JSONToString(freeWildBonusRateCount))

	fmt.Println("Normal RTP:", fmt.Sprintf("%.2f", RTP(normalScore, totalbetscore)), "normalScore:", normalScore, "TotalBet:", totalbetscore)
	fmt.Println("Scotter RTP:", fmt.Sprintf("%.2f", RTP(freeScore, totalbetscore)), "scotterScore:", freeScore, "TotalBet:", totalbetscore)
	fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	// fmt.Println("InRespintCount:", inRespintCount, "RespinCountMap:", foundation.JSONToString(respinCountMap))
	fmt.Println("------------------------------------------------------------------------------------------------")
	// }

}

func Test243ScotterGameRequest(test *testing.T) {

}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
