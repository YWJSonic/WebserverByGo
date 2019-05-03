package gamelogic

import (
	"fmt"
	"math/rand"
	"time"
)

type LogicResult struct {
	TotalWinMoney   int64
	PerTypeWinItem  [][]int
	PerTypeWinMoney [][]int
}

// GetGameResult All game result interface
func GetGameResult(gameid string, bet int, payrate []int, plate []int, scroll [][]int) interface{} {
	var result map[string]LogicResult

	normalGame := makeGamePlate(plate, scroll)

	fmt.Println(normalGame)
	return result

}

func makeGamePlate(plate []int, scroll [][]int) [][]int {
	var result [][]int
	var RmIndex int
	rand.Seed(time.Now().UnixNano())

	for i, value := range plate {
		RmIndex = rand.Intn(len(scroll[i]) - value)
		result = append(result, scroll[i][RmIndex:RmIndex+value])
	}

	return result
}
