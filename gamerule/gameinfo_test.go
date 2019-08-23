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
	var processLogSplit = 10

	for ii, iimax := 1, 6; ii < iimax; ii++ {
		var att []map[string]interface{}
		var result map[string]interface{}
		var otherdata map[string]interface{}
		var normalScore, freeScore, totalwinscore, totalbetscore, betindex, bonusscore, freegamecount, normalbonusganecount, freebonusganecount int64

		winLinePay = make(map[string]string)

		scotterRatecount := make(map[string]int64)
		Scotter1SymbolCount = make(map[string]int64)

		normalPayLineCount = make(map[string]int)
		normalGamePayLineScore := make(map[string]int64)
		normalgameScotter2SymbolCount = make(map[string]int64)

		freePayLineCount = make(map[string]int)
		freeGamePayLineScore := make(map[string]int64)
		freegameScotter2SymbolCount = make(map[string]int64)

		var playerid int64 = 24
		betmoney := GetBetMoney(betindex)
		attinfo := AttachInfo{PlayerID: playerid, Kind: GameIndex}
		att = attachInfoToAttachData(attinfo)
		RTPSetting = ii

		for index, max := 0, 100000000; index < max; index++ {
			totalbetscore += betmoney

			result, _, otherdata = GameRequest(playerid, betindex, att)
			if normalresult, ok := result["normalresult"]; ok {
				tmpbonusrate := ((normalresult).(map[string]interface{})["bonusrate"]).(int64)
				if tmpbonusrate > 0 {
					bonusscore += (tmpbonusrate * betmoney)
					totalwinscore += (tmpbonusrate * betmoney)
					normalbonusganecount++
					scotterRatecount[fmt.Sprintf("%v", tmpbonusrate)]++
				}

				gameresult, ok := ((normalresult).(map[string]interface{})["gameresult"]).([]interface{})
				if ok && len(gameresult) > 0 {
					for _, value := range gameresult {
						roundResult := value.(gameplate.InfoLine)
						score := roundResult.Score
						normalScore += score
						totalwinscore += score
						normalGamePayLineScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.Score
					}
				}
			}
			// free game
			if freeresult, ok := result["freeresult"]; ok {
				freegamecount++
				result := freeresult.([]map[string]interface{})

				for _, value := range result {
					tmpbonusrate := (value["bonusrate"]).(int64)
					if tmpbonusrate > 0 {
						freeScore += (tmpbonusrate * betmoney)
						totalwinscore += (tmpbonusrate * betmoney)
						freebonusganecount++
						scotterRatecount[fmt.Sprintf("%v", tmpbonusrate)]++
					}

					gameresult, ok := (value["gameresult"]).([]interface{})
					if ok && len(gameresult) > 0 {
						for _, value := range gameresult {
							roundResult := value.(gameplate.InfoLine)
							score := roundResult.Score
							freeScore += score
							totalwinscore += score
							freeGamePayLineScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.Score
						}
					}
				}
			}

			if otherdata["totalwinscore"].(int64) > 0 {
				// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
			}
			if index%(max/processLogSplit) == 0 {
				fmt.Println("Progress Rate:", float64(index)/float64(max)*100)
			}
		}

		fmt.Println("------------------------------------------------------------------------------------------------")
		fmt.Println("freeCount:", freegamecount)
		fmt.Println("normalbonusganecount:", normalbonusganecount)
		fmt.Println("freebonusganecount:", freebonusganecount)

		fmt.Println("winLinePay:", foundation.JSONToString(winLinePay))
		fmt.Println("scotterRatecount:", foundation.JSONToString(scotterRatecount))
		fmt.Println("Scotter1SymbolCount:", foundation.JSONToString(Scotter1SymbolCount))

		fmt.Println("normalPayLineCount:", foundation.JSONToString(normalPayLineCount))
		fmt.Println("normalGamePayLineScore:", foundation.JSONToString(normalGamePayLineScore))
		fmt.Println("normalgameScotter2SymbolCount:", foundation.JSONToString(normalgameScotter2SymbolCount))

		fmt.Println("freePayLineCount:", foundation.JSONToString(freePayLineCount))
		fmt.Println("freeGamePayLineScore:", foundation.JSONToString(freeGamePayLineScore))
		fmt.Println("freegameScotter2SymbolCount:", foundation.JSONToString(freegameScotter2SymbolCount))

		fmt.Println("Normal RTP:", fmt.Sprintf("%.2f", RTP(normalScore, totalbetscore)), "normalScore:", normalScore, "TotalBet:", totalbetscore)
		fmt.Println("freeScore RTP:", fmt.Sprintf("%.2f", RTP(freeScore, totalbetscore)), "scotterScore:", freeScore, "TotalBet:", totalbetscore)
		fmt.Println("Bonus RTP:", fmt.Sprintf("%.2f", RTP(bonusscore, totalbetscore)), "BonusScore:", bonusscore, "TotalBet:", totalbetscore)
		fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		fmt.Println("------------------------------------------------------------------------------------------------")
	}

}

func Test243ScotterGameRequest(test *testing.T) {

}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
