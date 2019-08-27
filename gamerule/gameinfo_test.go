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

	// var count int
	// quitFunc := make(chan int)
	quitChannel := make(chan int)
	analysisDataChannel := make(chan map[string]interface{})

	var emax = 6 //runtime.NumCPU()
	var rtpRange = 5
	for ii, iimax := 1, rtpRange; ii <= iimax; ii++ {

		RTPSetting = ii

		for e := 0; e < emax; e++ {
			go func() {
				analysisDataChannel <- mainLogic()
			}()
		}

		go analysis(emax, quitChannel, analysisDataChannel)
		<-quitChannel
	}
}

func analysis(emax int, quitChannel chan int, analysisDataChannel chan map[string]interface{}) {
	var Count int
	var normalScore, freeScore, totalwinscore, totalbetscore, bonusscore, freegamecount, normalbonusganecount, freebonusganecount int64
	var normalgameScotter2SymbolCount, scotterRatecount, scotter1SymbolCount map[string]int64
	var normalGamePayLineScore map[string]interface{}
	var freeGamePayLineScore, freegameScotter2SymbolCount map[string]int64

	normalGamePayLineScore = make(map[string]interface{})
	normalgameScotter2SymbolCount = make(map[string]int64)
	scotterRatecount = make(map[string]int64)
	scotter1SymbolCount = make(map[string]int64)
	freeGamePayLineScore = make(map[string]int64)
	freegameScotter2SymbolCount = make(map[string]int64)

	for {
		select {
		case data := <-analysisDataChannel:
			Count++
			normalScore += foundation.InterfaceToInt64(data["normalScore"])
			freeScore += foundation.InterfaceToInt64(data["freeScore"])
			totalwinscore += foundation.InterfaceToInt64(data["totalwinscore"])
			totalbetscore += foundation.InterfaceToInt64(data["totalbetscore"])
			bonusscore += foundation.InterfaceToInt64(data["bonusscore"])
			freegamecount += foundation.InterfaceToInt64(data["freegamecount"])
			normalbonusganecount += foundation.InterfaceToInt64(data["normalbonusganecount"])
			freebonusganecount += foundation.InterfaceToInt64(data["freebonusganecount"])

			normalgameScotter2SymbolCount = mapSum(normalgameScotter2SymbolCount, data["normalgameScotter2SymbolCount"].(map[string]int64))
			scotterRatecount = mapSum(scotterRatecount, data["scotterRatecount"].(map[string]int64))
			scotter1SymbolCount = mapSum(scotter1SymbolCount, data["scotter1SymbolCount"].(map[string]int64))

			freeGamePayLineScore = mapSum(freeGamePayLineScore, data["freeGamePayLineScore"].(map[string]int64))
			freegameScotter2SymbolCount = mapSum(freegameScotter2SymbolCount, data["freegameScotter2SymbolCount"].(map[string]int64))

			normalGamePayLineScore = foundation.AppendMap(normalGamePayLineScore, data["normalGamePayLineScore"].(map[string]interface{}))

			fmt.Println("Progress Rate:", float64(Count)/float64(emax)*100)
			if Count >= emax {
				fmt.Println("------------------------------------------------------------------------------------------------")

				fmt.Println("freeCount:", freegamecount)
				fmt.Println("normalbonusganecount:", normalbonusganecount)
				fmt.Println("freebonusganecount:", freebonusganecount)

				fmt.Println("scotterRatecount:", foundation.JSONToString(scotterRatecount))
				fmt.Println("scotter1SymbolCount:", foundation.JSONToString(scotter1SymbolCount))

				fmt.Println("normalGamePayLineScore:", foundation.JSONToString(normalGamePayLineScore))
				fmt.Println("normalgameScotter2SymbolCount:", foundation.JSONToString(normalgameScotter2SymbolCount))

				fmt.Println("freeGamePayLineScore:", foundation.JSONToString(freeGamePayLineScore))
				fmt.Println("freegameScotter2SymbolCount:", foundation.JSONToString(freegameScotter2SymbolCount))

				fmt.Println("Normal RTP:", fmt.Sprintf("%.2f", RTP(normalScore, totalbetscore)), "normalScore:", normalScore, "TotalBet:", totalbetscore)
				fmt.Println("freeScore RTP:", fmt.Sprintf("%.2f", RTP(freeScore, totalbetscore)), "scotterScore:", freeScore, "TotalBet:", totalbetscore)
				fmt.Println("Bonus RTP:", fmt.Sprintf("%.2f", RTP(bonusscore, totalbetscore)), "BonusScore:", bonusscore, "TotalBet:", totalbetscore)
				fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
				fmt.Println("------------------------------------------------------------------------------------------------")
				quitChannel <- 0
				return
			}
		}
	}
}

func mapSum(target, source map[string]int64) map[string]int64 {
	for key, value := range source {
		target[key] += value
	}
	return target
}

func mainLogic() map[string]interface{} {

	var processLogSplit = 10
	var playerid int64 = 24
	var betindex int64
	betmoney := GetBetMoney(betindex)
	var att []map[string]interface{}
	var result map[string]interface{}
	var otherdata map[string]interface{}
	analysisData := make(map[string]interface{})

	var normalScore, freeScore, totalwinscore, totalbetscore, bonusscore, freegamecount, normalbonusganecount, freebonusganecount int64
	normalgameScotter2SymbolCount := make(map[string]int64)
	scotterRatecount := make(map[string]int64)
	scotter1SymbolCount := make(map[string]int64)
	normalGamePayLineScore := make(map[string]interface{})
	freeGamePayLineScore := make(map[string]int64)
	freegameScotter2SymbolCount := make(map[string]int64)

	attinfo := AttachInfo{PlayerID: playerid, Kind: GameIndex}
	att = attachInfoToAttachData(attinfo)
	for index, max := 0, 2000000; index < max; index++ {

		totalbetscore += betmoney

		result, _, otherdata = GameRequest(playerid, betindex, att)
		if normalresult, ok := result["normalresult"]; ok {

			gameresult, ok := ((normalresult).(map[string]interface{})["gameresult"]).([]interface{})
			if ok && len(gameresult) > 0 {
				for _, value := range gameresult {
					roundResult := value.(gameplate.InfoLine)

					if roundResult.SpecialWinRate > 0 {
						if roundResult.LineWinRate == 2 {
							normalgameScotter2SymbolCount["3"]++
						} else if roundResult.LineWinRate == 5 {
							normalgameScotter2SymbolCount["4"]++
						} else if roundResult.LineWinRate == 10 {
							normalgameScotter2SymbolCount["5"]++
						}

						bonusscore += (roundResult.SpecialWinRate * betmoney)
						normalbonusganecount++
						scotterRatecount[fmt.Sprintf("%v", roundResult.SpecialWinRate)]++

					}

					normalScore += int64(roundResult.LineWinRate) * (betmoney / betLine)
					totalwinscore += roundResult.Score
					normalGamePayLineScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.Score

				}
			}
		}
		// free game
		if freeresult, ok := result["freeresult"]; ok {
			freegamecount++
			if len(result) == 10 {
				scotter1SymbolCount["3"]++
			} else if len(result) == 15 {
				scotter1SymbolCount["4"]++
			} else if len(result) == 20 {
				scotter1SymbolCount["5"]++
			}
			result := freeresult.([]map[string]interface{})

			for _, value := range result {
				gameresult, ok := (value["gameresult"]).([]interface{})
				if ok && len(gameresult) > 0 {
					for _, value := range gameresult {
						roundResult := value.(gameplate.InfoLine)

						if roundResult.SpecialWinRate > 0 {
							if roundResult.LineWinRate == 2 {
								freegameScotter2SymbolCount["3"]++
							} else if roundResult.LineWinRate == 5 {
								freegameScotter2SymbolCount["4"]++
							} else if roundResult.LineWinRate == 10 {
								freegameScotter2SymbolCount["5"]++
							}

							freebonusganecount++
							scotterRatecount[fmt.Sprintf("%v", roundResult.SpecialWinRate)]++

						}

						freeScore += roundResult.Score
						totalwinscore += roundResult.Score
						freeGamePayLineScore[fmt.Sprintf("%v", roundResult.LineSymbolNum)] = roundResult.Score
					}
				}
			}
		}

		if otherdata["totalwinscore"].(int64) > 0 {
			// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		}
		if index%(max/processLogSplit) == 0 {
			// fmt.Println("Progress Rate:", float64(index)/float64(max)*100)
		}
	}

	analysisData["totalbetscore"] = totalbetscore
	analysisData["normalScore"] = normalScore
	analysisData["freeScore"] = freeScore
	analysisData["totalwinscore"] = totalwinscore
	analysisData["totalbetscore"] = totalbetscore
	analysisData["bonusscore"] = bonusscore
	analysisData["freegamecount"] = freegamecount
	analysisData["normalbonusganecount"] = normalbonusganecount
	analysisData["freebonusganecount"] = freebonusganecount
	analysisData["normalgameScotter2SymbolCount"] = normalgameScotter2SymbolCount
	analysisData["scotterRatecount"] = scotterRatecount
	analysisData["scotter1SymbolCount"] = scotter1SymbolCount
	analysisData["normalGamePayLineScore"] = normalGamePayLineScore
	analysisData["freeGamePayLineScore"] = freeGamePayLineScore
	analysisData["freegameScotter2SymbolCount"] = freegameScotter2SymbolCount

	return analysisData
}

func Test243ScotterGameRequest(test *testing.T) {

}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
