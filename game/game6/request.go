package game6

import (
	"fmt"
	"math/rand"
	"time"

	"gitlab.com/WeberverByGo/foundation"
	"gitlab.com/WeberverByGo/game/gamesystem"
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

	if otherdata["isscatter2"].(int) == 1 {
		s2Result, otherdata, scatter2totalwin := scatter2Result(betMoney)
		totalWin += scatter2totalwin
		result["scatter2result"] = s2Result

		if otherdata["isscatter1"].(int) == 1 {
			s1Result, scatter1totalwin := scatter1Result(betMoney)
			totalWin += scatter1totalwin
			result["scatter1result"] = s1Result
		}

	} else if otherdata["isscatter1"].(int) == 1 {
		s1Result, scatter1totalwin := scatter1Result(betMoney)
		totalWin += scatter1totalwin
		result["scatter1result"] = s1Result
	}

	result["totalwinscore"] = totalWin

	if !(gamesystem.IsInTotalMoneyWinLimit(betMoney, totalWin) || gamesystem.IsInTotalBetRateWinLimit(betMoney, totalWin)) {
		return Result(betMoney, att...)
	}
	return result

}

var count int

func outputGame(betMoney int64) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	var result map[string]interface{}
	otherdata := make(map[string]interface{})
	islink := false

	ScrollIndex, plate := gamesystem.NewPlate(scrollSize, [][]int{Scroll1, Scroll2, Scroll3})

	count++
	plate = TestPlate(count % 4)
	gameresult := winresultArray(plate)
	fmt.Println(ScrollIndex, plate, gameresult)

	otherdata["isscatter1"] = 0
	otherdata["isscatter2"] = 0

	if isScatter1(plate) {
		otherdata["isscatter1"] = 1
	}
	if isScatter2(plate) {
		otherdata["isscatter2"] = 1
	}

	if len(gameresult) > 0 {
		totalScores = betMoney * int64(gameresult[0][3])
		islink = true
	}

	result = gamesystem.ResultMap(ScrollIndex, plate, totalScores, islink)
	return result, otherdata, foundation.InterfaceToInt64(result["scores"])
}

// GameResult ...
func winresultArray(plate []int) [][]int {
	var result [][]int

	for _, ItemResult := range itemResults {
		if isWin(plate, ItemResult) {
			result = append(result, ItemResult)
			if isSingleLine {
				break
			}
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
			if plate == scatter1 || plate == scatter2 {
				IsWin = false
			} else {
				IsWin = true
			}
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

func isScatter1(itemResult []int) bool {
	if itemResult[3] == -100 {
		return true
	}
	return false
}
func isScatter2(itemResult []int) bool {
	if itemResult[3] == -101 {
		return true
	}
	return false
}

func isAny7(item int) bool {
	if item == 4 || item == 5 || item == 6 {
		return true
	}
	return false
}

func isAnyBay(item int) bool {
	if item == 7 || item == 8 || item == 9 || item == 10 {
		return true
	}
	return false
}

func scatter1Result(betMoney int64) (map[string]interface{}, int64) {
	var result map[string]interface{}

	scatterIndex := foundation.RangeRandom(Scatter1Range[Scatter1Setting])
	scatterWinRate := Scatter1WinRate[scatterIndex]
	totalScores := betMoney * scatterWinRate
	result = gamesystem.ResultMap([]int{scatterIndex}, []int{scatterIndex}, totalScores, true)
	return result, totalScores
}
func scatter2Result(betMoney int64) (map[string]interface{}, map[string]interface{}, int64) {
	var result map[string]interface{}
	otherdata := make(map[string]interface{})
	var totalScores int64

	scatterIndex := foundation.RangeRandom(Scatter2Range[Scatter2Setting])
	scatterWinRate := Scatter2WinRate[scatterIndex]

	otherdata["isscatter1"] = 0

	if scatterWinRate < 0 {
		otherdata["isscatter1"] = 1
		totalScores = betMoney * scatterWinRate
	}
	result = gamesystem.ResultMap([]int{scatterIndex}, []int{scatterIndex}, totalScores, true)
	return result, otherdata, totalScores

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
	demoplate := [][]int{{0, 0, 0}, {1, 1, 1}, {3, 2, 9}, {3, 8, 2}, {2, 5, 7}}
	plate := demoplate[index]
	return plate
}
