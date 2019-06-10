package game6

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/WeberverByGo/foundation"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// BetRate ...
func BetRate() []int64 {
	return betRate
}

// Scroll ...
func Scroll() interface{} {
	scrollmap := map[string][][]int{
		"normalreel": {Scroll1, Scroll2, Scroll3},
	}
	return scrollmap
}

// Result no att value
func Result(betMoney int64, att ...interface{}) map[string]interface{} {
	var result = make(map[string]interface{})
	var totalWin int64

	normalresult, otherdata, normaltotalwin := outputGame(betMoney)
	result = otherdata
	result["normalresult"] = normalresult
	totalWin += normaltotalwin

	result["totalwinscore"] = totalWin

	if (totalWin / betMoney) > limitWinBetRate {
		return Result(betMoney, att...)
	}

	return result

}

var count int

func outputGame(betMoney int64) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	otherdata := make(map[string]interface{})
	result := make(map[string]interface{})

	ScrollIndex, plate := newPlate(scrollSize, [][]int{Scroll1, Scroll2, Scroll3})

	count++
	plate = TestPlate(count % 4)
	gameresult := winresultArray(plate)
	fmt.Println(ScrollIndex, plate, gameresult)

	result["plateindex"] = ScrollIndex
	result["plate"] = plate
	result["scores"] = 0
	result["islink"] = 0

	if len(gameresult) > 0 {
		totalScores = betMoney * int64(gameresult[0][3])
		result["scores"] = totalScores
		result["islink"] = 1
	}

	return result, otherdata, foundation.InterfaceToInt64(result["scores"])
}

// NewPlate ...
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

// GameResult ...
func winresultArray(plate []int) [][]int {
	var result [][]int

	for _, ItemResult := range ItemResults {
		if isWin(plate, ItemResult) {
			result = append(result, ItemResult)
		}
	}
	return result
}

func isWin(plates []int, result []int) bool {
	IsWin := false

	for i, plate := range plates {
		IsWin = false
		switch result[i] {
		case plate:
			IsWin = true
		case -1000:
			IsWin = true
		case -1001:
			if isAny7(plate) {
				IsWin = true
			}
		case -1002:
			if isAnyBay(plate) {
				IsWin = true
			}
		}
		if !IsWin {
			return IsWin
		}
	}

	return IsWin
}

func isAny7(Item int) bool {
	if Item == 4 || Item == 5 || Item == 6 || Item == 7 {
		return true
	}
	return false
}

func isAnyBay(Item int) bool {
	if Item == 7 || Item == 8 || Item == 9 || Item == 10 {
		return true
	}
	return false
}

func scatter1() (int, int) {

	scatterIndex := foundation.RangeRandom(Scatter1Range[Scatter1Setting])
	scatterBet := Scatter1Score[scatterIndex]
	return scatterIndex, scatterBet
}
func scatter2() (int, int) {

	scatterIndex := foundation.RangeRandom(Scatter2Range[Scatter2Setting])
	scatterBet := Scatter2Score[scatterIndex]
	return scatterIndex, scatterBet
}

// TestPlate ...
func TestPlate(index int) []int {
	// var vcount int
	// plate := []int{0, 0, 0}
	// demoplate := [][]int{{0, 0, 0}, {0, 0, 5}, {0, 0, 6}, {0, 0, 7}, {0, 0, 8}, {0, 0, 9}, {0, 0, 10}, {0, 1, 0}, {0, 1, 5}, {0, 1, 6}, {0, 1, 7}, {0, 1, 8}, {0, 1, 9}, {0, 1, 10}, {0, 5, 0}, {0, 5, 5}, {0, 5, 6}, {0, 5, 7}, {0, 5, 8}, {0, 5, 9}, {0, 5, 10}, {0, 6, 0}, {0, 6, 5}, {0, 6, 6}, {0, 6, 7}, {0, 6, 8}, {0, 6, 9}, {0, 6, 10}, {0, 7, 0}, {0, 7, 5}, {0, 7, 6}, {0, 7, 7}, {0, 7, 8}, {0, 7, 9}, {0, 7, 10}, {0, 8, 0}, {0, 8, 5}, {0, 8, 6}, {0, 8, 7}, {0, 8, 8}, {0, 8, 9}, {0, 8, 10}, {0, 9, 0}, {0, 9, 5}, {0, 9, 6}, {0, 9, 7}, {0, 9, 8}, {0, 9, 9}, {0, 9, 10}, {0, 10, 0}, {0, 10, 5}, {0, 10, 6}, {0, 10, 7}, {0, 10, 8}, {0, 10, 9}, {0, 10, 10}, {5, 0, 0}, {5, 0, 5}, {5, 0, 6}, {5, 0, 7}, {5, 0, 8}, {5, 0, 9}, {5, 0, 10}, {5, 1, 0}, {5, 1, 5}, {5, 1, 6}, {5, 1, 7}, {5, 1, 8}, {5, 1, 9}, {5, 1, 10}, {5, 5, 0}, {5, 5, 5}, {5, 5, 6}, {5, 5, 7}, {5, 5, 8}, {5, 5, 9}, {5, 5, 10}, {5, 6, 0}, {5, 6, 5}, {5, 6, 6}, {5, 6, 7}, {5, 6, 8}, {5, 6, 9}, {5, 6, 10}, {5, 7, 0}, {5, 7, 5}, {5, 7, 6}, {5, 7, 7}, {5, 7, 8}, {5, 7, 9}, {5, 7, 10}, {5, 8, 0}, {5, 8, 5}, {5, 8, 6}, {5, 8, 7}, {5, 8, 8}, {5, 8, 9}, {5, 8, 10}, {5, 9, 0}, {5, 9, 5}, {5, 9, 6}, {5, 9, 7}, {5, 9, 8}, {5, 9, 9}, {5, 9, 10}, {5, 10, 0}, {5, 10, 5}, {5, 10, 6}, {5, 10, 7}, {5, 10, 8}, {5, 10, 9}, {5, 10, 10}, {6, 0, 0}, {6, 0, 5}, {6, 0, 6}, {6, 0, 7}, {6, 0, 8}, {6, 0, 9}, {6, 0, 10}, {6, 1, 0}, {6, 1, 5}, {6, 1, 6}, {6, 1, 7}, {6, 1, 8}, {6, 1, 9}, {6, 1, 10}, {6, 5, 0}, {6, 5, 5}, {6, 5, 6}, {6, 5, 7}, {6, 5, 8}, {6, 5, 9}, {6, 5, 10}, {6, 6, 0}, {6, 6, 5}, {6, 6, 6}, {6, 6, 7}, {6, 6, 8}, {6, 6, 9}, {6, 6, 10}, {6, 7, 0}, {6, 7, 5}, {6, 7, 6}, {6, 7, 7}, {6, 7, 8}, {6, 7, 9}, {6, 7, 10}, {6, 8, 0}, {6, 8, 5}, {6, 8, 6}, {6, 8, 7}, {6, 8, 8}, {6, 8, 9}, {6, 8, 10}, {6, 9, 0}, {6, 9, 5}, {6, 9, 6}, {6, 9, 7}, {6, 9, 8}, {6, 9, 9}, {6, 9, 10}, {6, 10, 0}, {6, 10, 5}, {6, 10, 6}, {6, 10, 7}, {6, 10, 8}, {6, 10, 9}, {6, 10, 10}, {7, 0, 0}, {7, 0, 5}, {7, 0, 6}, {7, 0, 7}, {7, 0, 8}, {7, 0, 9}, {7, 0, 10}, {7, 1, 0}, {7, 1, 5}, {7, 1, 6}, {7, 1, 7}, {7, 1, 8}, {7, 1, 9}, {7, 1, 10}, {7, 5, 0}, {7, 5, 5}, {7, 5, 6}, {7, 5, 7}, {7, 5, 8}, {7, 5, 9}, {7, 5, 10}, {7, 6, 0}, {7, 6, 5}, {7, 6, 6}, {7, 6, 7}, {7, 6, 8}, {7, 6, 9}, {7, 6, 10}, {7, 7, 0}, {7, 7, 5}, {7, 7, 6}, {7, 7, 7}, {7, 7, 8}, {7, 7, 9}, {7, 7, 10}, {7, 8, 0}, {7, 8, 5}, {7, 8, 6}, {7, 8, 7}, {7, 8, 8}, {7, 8, 9}, {7, 8, 10}, {7, 9, 0}, {7, 9, 5}, {7, 9, 6}, {7, 9, 7}, {7, 9, 8}, {7, 9, 9}, {7, 9, 10}, {7, 10, 0}, {7, 10, 5}, {7, 10, 6}, {7, 10, 7}, {7, 10, 8}, {7, 10, 9}, {7, 10, 10}, {8, 0, 0}, {8, 0, 5}, {8, 0, 6}, {8, 0, 7}, {8, 0, 8}, {8, 0, 9}, {8, 0, 10}, {8, 1, 0}, {8, 1, 5}, {8, 1, 6}, {8, 1, 7}, {8, 1, 8}, {8, 1, 9}, {8, 1, 10}, {8, 5, 0}, {8, 5, 5}, {8, 5, 6}, {8, 5, 7}, {8, 5, 8}, {8, 5, 9}, {8, 5, 10}, {8, 6, 0}, {8, 6, 5}, {8, 6, 6}, {8, 6, 7}, {8, 6, 8}, {8, 6, 9}, {8, 6, 10}, {8, 7, 0}, {8, 7, 5}, {8, 7, 6}, {8, 7, 7}, {8, 7, 8}, {8, 7, 9}, {8, 7, 10}, {8, 8, 0}, {8, 8, 5}, {8, 8, 6}, {8, 8, 7}, {8, 8, 8}, {8, 8, 9}, {8, 8, 10}, {8, 9, 0}, {8, 9, 5}, {8, 9, 6}, {8, 9, 7}, {8, 9, 8}, {8, 9, 9}, {8, 9, 10}, {8, 10, 0}, {8, 10, 5}, {8, 10, 6}, {8, 10, 7}, {8, 10, 8}, {8, 10, 9}, {8, 10, 10}, {9, 0, 0}, {9, 0, 5}, {9, 0, 6}, {9, 0, 7}, {9, 0, 8}, {9, 0, 9}, {9, 0, 10}, {9, 1, 0}, {9, 1, 5}, {9, 1, 6}, {9, 1, 7}, {9, 1, 8}, {9, 1, 9}, {9, 1, 10}, {9, 5, 0}, {9, 5, 5}, {9, 5, 6}, {9, 5, 7}, {9, 5, 8}, {9, 5, 9}, {9, 5, 10}, {9, 6, 0}, {9, 6, 5}, {9, 6, 6}, {9, 6, 7}, {9, 6, 8}, {9, 6, 9}, {9, 6, 10}, {9, 7, 0}, {9, 7, 5}, {9, 7, 6}, {9, 7, 7}, {9, 7, 8}, {9, 7, 9}, {9, 7, 10}, {9, 8, 0}, {9, 8, 5}, {9, 8, 6}, {9, 8, 7}, {9, 8, 8}, {9, 8, 9}, {9, 8, 10}, {9, 9, 0}, {9, 9, 5}, {9, 9, 6}, {9, 9, 7}, {9, 9, 8}, {9, 9, 9}, {9, 9, 10}, {9, 10, 0}, {9, 10, 5}, {9, 10, 6}, {9, 10, 7}, {9, 10, 8}, {9, 10, 9}, {9, 10, 10}, {10, 0, 0}, {10, 0, 5}, {10, 0, 6}, {10, 0, 7}, {10, 0, 8}, {10, 0, 9}, {10, 0, 10}, {10, 1, 0}, {10, 1, 5}, {10, 1, 6}, {10, 1, 7}, {10, 1, 8}, {10, 1, 9}, {10, 1, 10}, {10, 5, 0}, {10, 5, 5}, {10, 5, 6}, {10, 5, 7}, {10, 5, 8}, {10, 5, 9}, {10, 5, 10}, {10, 6, 0}, {10, 6, 5}, {10, 6, 6}, {10, 6, 7}, {10, 6, 8}, {10, 6, 9}, {10, 6, 10}, {10, 7, 0}, {10, 7, 5}, {10, 7, 6}, {10, 7, 7}, {10, 7, 8}, {10, 7, 9}, {10, 7, 10}, {10, 8, 0}, {10, 8, 5}, {10, 8, 6}, {10, 8, 7}, {10, 8, 8}, {10, 8, 9}, {10, 8, 10}, {10, 9, 0}, {10, 9, 5}, {10, 9, 6}, {10, 9, 7}, {10, 9, 8}, {10, 9, 9}, {10, 9, 10}, {10, 10, 0}, {10, 10, 5}, {10, 10, 6}, {10, 10, 7}, {10, 10, 8}, {10, 10, 9}, {10, 10, 10}}
	// plate = demoplate[index]
	// for i, imax := 0, 10; i <= imax; i++ {
	// 	for j, jmax := 0, 10; j <= jmax; j++ {
	// 		for k, kmax := 0, 10; k <= kmax; k++ {
	// 			plate[0] = i
	// 			plate[1] = j
	// 			plate[2] = k

	// 			fmt.Println(vcount, ":", plate, winresultArray(plate))
	// 			if index == vcount {
	// 				return plate
	// 			}
	// 			vcount++
	// 		}
	// 	}
	// }

	// for _, plate := range demoplate {
	// 	fmt.Println(vcount, ":", plate, winresultArray(plate))
	// }
	demoplate := [][]int{{7, 1, 0}, {8, 1, 9}, {6, 9, 5}, {10, 1, 0}}
	plate := demoplate[index]
	return plate
}
