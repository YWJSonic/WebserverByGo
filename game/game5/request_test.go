package game5_test

import (
	"fmt"
	"testing"

	"gitlab.com/WeberverByGo/game/game5"
)

func TestResult(t *testing.T) {
	var totalwin, totalbet int64 = 0, 0
	var freecount = 0
	var bet int64 = 100
	var result map[string]interface{}
	for index3 := 1; index3 < 4; index3++ {
		game5.RespinSetting = index3
		for index2 := 0; index2 < 5; index2++ {
			totalwin, totalbet, freecount = 0, 0, 0
			for index, max := 0, 10000000; index < max; index++ {
				totalbet += bet
				result = game5.Result(bet, freecount)
				freecount = result["freecount"].(int)
				totalwin += result["totalwinscore"].(int64)
			}
			fmt.Printf("RTP:%f\n  totalwin:%d \n totalbet:%d \n freecount:%d \n", float64(totalwin)/float64(totalbet)*100, totalwin, totalbet, freecount)
		}
	}
}
