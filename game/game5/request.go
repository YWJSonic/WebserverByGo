package game5

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/math"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GameIndex ...
func GameIndex() int64 {
	return gameIndex
}

// BetRate ...
func BetRate() interface{} {
	return betRate
}

// Scroll ...
func Scroll() interface{} {
	scrollmap := map[string][][]int{
		"normalscroll": normalScroll,
		"freescroll":   freeScroll,
		"respinscroll": {respuinScroll()},
	}
	return scrollmap
}

// Result ...
func Result(BetIndex int64, FreeCount int) interface{} {
	var result = make(map[string]interface{})
	var totalWin int64
	Bet := betRate[BetIndex]

	fmt.Println("----")
	normalresult, normaltotalwin := outputGame(Bet, FreeCount)

	totalWin += normaltotalwin
	result["normalresult"] = normalresult

	if normalresult["isfreegame"].(int) == 1 {
		fmt.Println("----")
		freeresult, freetotalwin := outputFreeSpin()
		totalWin += freetotalwin
		result["freegame"] = freeresult

	}
	if normalresult["isrespin"].(int) == 1 {
		fmt.Println("----")
		respinresult, respintotalwin := outRespin()
		totalWin += respintotalwin
		result["respin"] = respinresult

	}

	result["totalwinscore"] = totalWin
	return result
}

func outputGame(bet int64, freecount int) (map[string]interface{}, int64) {
	var totalScores int64
	result := make(map[string]interface{})

	ScrollIndex, plate := newPlate(scrollSize, normalScroll)
	gameresult := winresultArray(plate)
	fmt.Println(ScrollIndex, plate, gameresult)

	result["plateindex"] = ScrollIndex
	result["plate"] = plate
	result["isfreegame"] = 0
	result["freecount"] = freecount
	result["isrespin"] = 0
	result["scores"] = 0

	if isFreeGameCount(plate) {
		freecount++
		if freecount >= freeGameTrigger {
			result["isfreegame"] = 1
			result["freecount"] = 0
		} else {
			result["freecount"] = freecount
		}
	}

	if isRespin(plate) {
		result["isrespin"] = 1
	}

	if len(gameresult) > 0 {
		totalScores = bet * int64(gameresult[0][3])
		result["scores"] = totalScores
	}

	return result, foundation.InterfaceToInt64(result["scores"])
}

func outputFreeSpin() ([]interface{}, int64) {
	var result []interface{}
	var totalScores int64

	for i, max := 0, 5; i < max; i++ {
		freeresult := make(map[string]interface{})
		ScrollIndex, plate := newPlate(scrollSize, freeScroll)
		gr := winresultArray(plate)
		fmt.Println(ScrollIndex, plate, gr)

		freeresult["plateindex"] = ScrollIndex
		freeresult["plate"] = plate
		freeresult["scores"] = 0
		if len(gr) > 0 {
			freeresult["scores"] = gr[0][2]
			totalScores += foundation.InterfaceToInt64(freeresult["scores"])
		}

		result = append(result, freeresult)
	}
	return result, totalScores
}

func outRespin() (map[string]interface{}, int64) {
	var totalscores int64

	ScrollIndex, plate := newPlate([]int{1}, [][]int{respuinScroll()})
	gr := respinResult(plate)
	fmt.Println(ScrollIndex, plate, gr)

	result := make(map[string]interface{})
	result["plateindex"] = ScrollIndex
	result["plate"] = plate
	result["scores"] = 0
	if len(gr) > 0 {
		result["scores"] = gr[0][1]
		totalscores += foundation.InterfaceToInt64(result["scores"])
	}

	return result, totalscores
}

// RespuinScroll ...
func respuinScroll() []int {
	if RespinSetting == 1 {
		return respinScroll1
	} else if RespinSetting == 2 {
		return respinScroll2
	} else {
		return respinScroll3
	}
}

// winresultArray ...
func winresultArray(plate []int) [][]int {
	var result [][]int
	var dynamicresult []int

	for _, ItemResult := range itemResults {
		if isWin(plate, ItemResult) {

			if isDynamicResult(ItemResult) {
				dynamicresult = dynamicScore(plate, ItemResult)
				result = append(result, dynamicresult)
				if isSingleLine {
					break
				}
			} else {
				result = append(result, ItemResult)
				if isSingleLine {
					break
				}
			}
		}
	}

	return result
}

// RespinResult result 0: icon index, 1: win rate
func respinResult(plate []int) [][]int {
	var result [][]int

	switch plate[0] {
	case 2:
		result = append(result, []int{2, 5})
	case 3:
		result = append(result, []int{3, 7})
	case 4:
		result = append(result, []int{4, 10})
	}

	return result
}

// IsFreeGameCount ...
func isFreeGameCount(plate []int) bool {
	if plate[1] == 1 {
		return true
	}
	return false

}

// IsRespin ...
func isRespin(plate []int) bool {
	if plate[0] != 10 && plate[1] == 1 && plate[2] == 0 {
		return true
	}
	return false

}

func newPlate(plateSize []int, scroll [][]int) ([]int, []int) {
	var ScrollIndex []int
	var plate []int
	var index int

	for i := range plateSize {
		index = rand.Intn(len(scroll[i]))
		plate = append(plate, scroll[i][index])
		ScrollIndex = append(ScrollIndex, index)
	}

	return ScrollIndex, plate
}

func isWin(plates []int, result []int) bool {
	IsWin := false

	for i, plate := range plates {
		IsWin = false

		if plate == space {
			return IsWin
		}

		switch result[i] {
		case plate:
			IsWin = true
		case -1000:
			IsWin = true
		case -1001: // any 7
			if math.IsInclude(plate, []int{5, 6}) {
				IsWin = true
			}
		case -1002: // any bar
			if math.IsInclude(plate, []int{6, 7, 8, 9}) {
				IsWin = true
			}
		}
		if !IsWin {
			return IsWin
		}
	}

	return IsWin
}

func dynamicScore(plant, currendResult []int) []int {
	var dynamicresult []int

	switch currendResult[3] {
	case -100:
		dynamicresult = currendResult[:3]
		for _, result := range itemResults {
			if result[0] == plant[0] {
				dynamicresult = append(dynamicresult, result[3])
				break
			}
		}
	}

	return dynamicresult
}

func isDynamicResult(result []int) bool {
	if result[3] < 0 {
		return true
	}
	return false
}

// TestPlate ...
func TestPlate(index int) []int {
	// var vcount int
	plate := []int{0, 0, 0}
	demoplate := [][]int{
		{0, 0, 0},
		{5, 5, 5},
		{6, 6, 6},
		{7, 7, 7},
		{8, 8, 8},
		{9, 9, 9},
		{0, 1, 0},
		{5, 6, 8},
	}
	plate = demoplate[index]
	// for i, imax := 0, 10; i <= imax; i++ {
	// 	for j, jmax := 0, 10; j <= jmax; j++ {
	// 		for k, kmax := 0, 10; k <= kmax; k++ {
	// 			plate[0] = i
	// 			plate[1] = j
	// 			plate[2] = k

	// 			fmt.Println(vcount, ":", plate, GameResult(plate))
	// 			if index == vcount {
	// 				return plate
	// 			}
	// 			vcount++
	// 		}
	// 	}
	// }
	return plate
}
