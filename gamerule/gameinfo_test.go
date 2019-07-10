package gamerule

import (
	"fmt"
	"testing"

	"gitlab.com/ServerUtility/foundation"
	moneypool "gitlab.com/WeberverByGo/servercontrol"
)

// TestGameRequest ...
func TestGameRequest(test *testing.T) {
	var att []map[string]interface{}
	var newatt []map[string]interface{}
	var result map[string]interface{}
	var otherdata map[string]int64
	var totalwinscore, totalbetscore, betindex int64
	var inRespintCount int64
	var logplate string
	respinCountMap := make(map[int]int64)
	normalPlayteCount := make(map[string]int64)
	respinPlayteCount := make(map[string]int64)
	normalPlayteWin := make(map[string]int64)
	respinPlayteWin := make(map[string]int64)
	betmoney := GetBetMoney(betindex)
	attinfo := AttachInfo{PlayerID: 4, Kind: 7, JackPartBonusPoolx2: 0, JackPartBonusPoolx3: 0, JackPartBonusPoolx5: 0}
	att = attachInfoToAttachData(attinfo)

	for index, max := 0, 100000000; index < max; index++ {
		totalbetscore += betmoney

		result, newatt, otherdata = GameRequest(4, 0, att)
		att = newatt
		if normalresult, ok := result["normalresult"]; ok {
			logplate = fmt.Sprintf("%v", (normalresult).(map[string]interface{})["plate"])
			normalPlayteCount[logplate]++
			if _, ok := normalPlayteWin[logplate]; !ok {
				normalPlayteWin[logplate] = (normalresult.(map[string]interface{})["scores"]).(int64)
			}
		}

		if respinsresult, ok := result["respin"]; ok {
			inRespintCount++
			respinCountMap[len(respinsresult.([]interface{}))]++
			for _, respinresult := range respinsresult.([]interface{}) {
				logplate = fmt.Sprintf("%v", (respinresult).(map[string]interface{})["plate"])
				respinPlayteCount[logplate]++

				if _, ok := respinPlayteWin[logplate]; !ok {
					respinPlayteWin[logplate] = (respinresult.(map[string]interface{})["scores"]).(int64)
				}
			}
		}

		moneypool.RTPControl(betmoney, otherdata["totalwinscore"])
		totalwinscore += otherdata["totalwinscore"]
		if otherdata["totalwinscore"] > 0 {
			// fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
		}
	}

	fmt.Println("NormalPlayteCount:", foundation.JSONToString(normalPlayteCount))
	fmt.Println("RespinPlayteCount:", foundation.JSONToString(respinPlayteCount))
	fmt.Println("RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	fmt.Println("InRespintCount:", inRespintCount, "RespinCountMap:", foundation.JSONToString(respinCountMap))
	// fmt.Println(foundation.JSONToString(normalPlayteWin))
	// fmt.Println(foundation.JSONToString(respinPlayteWin))
}

func TestNormalGameRTP(test *testing.T) {
	var BetIndex int64
	var totalwinscore, totalbetscore int64
	betMoney := GetBetMoney(BetIndex)
	attinfo := AttachInfo{PlayerID: 4, Kind: 7, JackPartBonusPoolx2: 0, JackPartBonusPoolx3: 0, JackPartBonusPoolx5: 0}

	for index, max := 0, 10000000; index < max; index++ {
		totalbetscore += betMoney

		_, _, normaltotalwin := outputGame(betMoney, &attinfo)
		totalwinscore += normaltotalwin
		fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	}
}

func TestRespinGameRTP(test *testing.T) {
	var BetIndex int64
	var totalwinscore, totalbetscore int64
	betMoney := GetBetMoney(BetIndex)
	attinfo := AttachInfo{PlayerID: 4, Kind: 7, JackPartBonusPoolx2: 0, JackPartBonusPoolx3: 0, JackPartBonusPoolx5: 0}

	for index, max := 0, 100000000; index < max; index++ {
		totalbetscore += betMoney

		_, respintotalwin := outRespin(betMoney, &attinfo)
		totalwinscore += respintotalwin
		fmt.Println("Index:", index, "RTP:", fmt.Sprintf("%.2f", RTP(totalwinscore, totalbetscore)), "TotalWin:", totalwinscore, "TotalBet:", totalbetscore)
	}
}
func TestNormalResultArray(test *testing.T) {
	var result [][]int
	scroll := [][]int{
		{0, 4, 5, 6, 7, 8},
		{1, 2, 3, 4, 5, 6, 7, 8},
		{0, 4, 5, 6, 7, 8},
	}
	for _, scroll1 := range scroll[0] {
		for _, scroll2 := range scroll[1] {
			for _, scroll3 := range scroll[2] {
				result = normalResultArray(scroll)
				if len(result) > 0 {
					fmt.Println("plate1:", scroll1, "plate2:", scroll2, "plate3:", scroll3, "Result:", result)
				}
			}
		}
	}
}

func TestRespinResultArray(test *testing.T) {
	var result [][]int
	scroll := [][]int{
		{0},
		{1, 2, 3, 4, 5, 6, 7, 8},
		{0},
	}
	for _, scroll1 := range scroll[0] {
		for _, scroll2 := range scroll[1] {
			for _, scroll3 := range scroll[2] {
				result = respinResultArray([]int{scroll1, scroll2, scroll3})
				if len(result) > 0 {
					fmt.Println("plate1:", scroll1, "plate2:", scroll2, "plate3:", scroll3, "Result:", result)
				}
			}
		}
	}
}

func RTP(totalwin, TotalBet int64) float64 {

	return float64(totalwin) / float64(TotalBet) * 100
}
