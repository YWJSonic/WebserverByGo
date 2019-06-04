package game

import (
	"math/rand"
	"time"

	"gitlab.com/WeberverByGo/foundation"
)

// 花開富貴

// GameIndex game sort id
const GameIndex = 6

// Scroll1 roll 1
var Scroll1 = []int{5, 10, 2, 9, 7, 4, 1, 1, 1, 1, 3, 6, 9, 0, 0, 0, 8, 5, 7, 4, 8, 9, 5, 9, 7, 5, 6, 1, 1, 1, 1, 7, 6, 7, 8, 9, 4, 1, 1, 1, 1, 8, 7, 10, 1, 0, 1, 10, 9, 4, 10, 5, 10, 8, 7, 8, 10, 8, 7, 6, 1, 0, 1, 8, 6, 10, 1, 0, 1, 10, 4, 6, 10, 4, 8, 10}

// Scroll2 roll 2
var Scroll2 = []int{4, 5, 8, 2, 5, 10, 4, 3, 7, 9, 0, 0, 0, 10, 8, 4, 6, 5, 9, 1, 0, 1, 8, 6, 4, 9, 1, 0, 1, 9, 10, 4, 7, 8, 5, 6, 10, 4, 1, 1, 1, 1, 6, 7, 1, 0, 1, 7, 6, 8, 10, 9, 8, 1, 0, 1, 9, 7, 9, 8, 3, 5, 7, 9, 4, 5, 6, 10, 8, 9, 1, 0, 1, 4, 10, 8}

// Scroll3 roll 3
var Scroll3 = []int{5, 8, 7, 5, 10, 4, 5, 10, 8, 5, 2, 6, 0, 0, 0, 10, 9, 4, 7, 5, 4, 3, 8, 1, 1, 1, 1, 6, 10, 1, 0, 1, 9, 6, 8, 7, 8, 1, 0, 1, 9, 4, 7, 6, 10, 1, 0, 1, 10, 9, 8, 4, 6, 8, 9, 6, 1, 0, 1, 4, 9, 8, 0, 0, 0, 6, 7, 8, 4, 9, 6, 7, 9, 4, 5, 9}

// Scatter1RangeSum ...
var Scatter1RangeSum = 50 + 50 + 30 + 20 + 10 + 5

// Scatter1Setting Scatter1Range Index
var Scatter1Setting = 3

// Scatter1Range ...
var Scatter1Range = [][]int{
	[]int{80, 70, 15, 8, 3, 1},
	[]int{80, 70, 15, 8, 3, 1},
	[]int{80, 70, 20, 15, 5, 5},
	[]int{50, 50, 30, 20, 10, 5}}

// Scatter1Score bonus game 1 rate sum
var Scatter1Score = []int{10, 20, 80, 200, 500, 1000}

// Scatter2RangeSum ...
var Scatter2RangeSum = 50 + 50 + 50 + 30 + 20 + 21

// Scatter2Setting Scatter2Range Index
var Scatter2Setting = 3

// Scatter2Range bonus game 2 rate sum
var Scatter2Range = [][]int{
	[]int{70, 70, 50, 13, 6, 15},
	[]int{70, 70, 50, 30, 20, 21},
	[]int{70, 70, 50, 30, 20, 21},
	[]int{50, 50, 50, 30, 20, 21}}

// Scatter2Score ...
var Scatter2Score = []int{3, 8, 20, 50, 100, 0}

var items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// ItemResults 0~10 item
// -1000 any
// -1001 any 7
// -1002 any bar
// -100 bonus game 1
// -101 bonus game 2
var ItemResults = [][]int{
	{0, 0, 0, -100},
	{1, 1, 1, -101},
	{2, 2, 3, 1000},
	{2, 3, 2, 1000},
	{3, 2, 2, 1000},
	{2, 3, 3, 900},
	{3, 2, 3, 900},
	{3, 3, 2, 900},
	{2, 2, 2, 500},
	{3, 3, 3, 300},
	{2, 3, -1000, 100},
	{2, -1000, 3, 100},
	{-1000, 2, 3, 100},
	{3, 2, -1000, 100},
	{3, -1000, 2, 100},
	{-1000, 3, 2, 100},
	{2, 4, 4, 60},
	{2, 5, 5, 45},
	{2, 6, 6, 30},
	{2, 7, 7, 24},
	{2, 8, 8, 15},
	{2, 9, 9, 9},
	{2, 10, 10, 6},
	{4, 2, 4, 60},
	{5, 2, 5, 45},
	{6, 2, 6, 30},
	{7, 2, 7, 24},
	{8, 2, 8, 15},
	{9, 2, 9, 9},
	{10, 2, 10, 6},
	{4, 4, 2, 60},
	{5, 5, 2, 45},
	{6, 6, 2, 30},
	{7, 7, 2, 24},
	{8, 8, 2, 15},
	{9, 9, 2, 9},
	{10, 10, 2, 6},
	{2, -1000, -1000, 3},
	{-1000, 2, -1000, 3},
	{-1000, -1000, 2, 3},
	{4, 2, 2, 60},
	{5, 2, 2, 45},
	{6, 2, 2, 30},
	{7, 2, 2, 24},
	{8, 2, 2, 15},
	{9, 2, 2, 9},
	{10, 2, 2, 6},
	{0, 2, 2, 3},
	{1, 2, 2, 3},
	{2, 4, 2, 60},
	{2, 5, 2, 45},
	{2, 6, 2, 30},
	{2, 7, 2, 24},
	{2, 8, 2, 15},
	{2, 9, 2, 9},
	{2, 10, 2, 6},
	{2, 0, 2, 3},
	{2, 1, 2, 3},
	{2, 2, 4, 60},
	{2, 2, 5, 45},
	{2, 2, 6, 30},
	{2, 2, 7, 24},
	{2, 2, 8, 15},
	{2, 2, 9, 9},
	{2, 2, 10, 6},
	{2, 2, 0, 3},
	{2, 2, 1, 3},
	{3, 4, 4, 40},
	{3, 5, 5, 30},
	{3, 6, 6, 20},
	{3, 7, 7, 16},
	{3, 8, 8, 10},
	{3, 9, 9, 6},
	{3, 10, 10, 4},
	{4, 3, 4, 40},
	{5, 3, 5, 30},
	{6, 3, 6, 20},
	{7, 3, 7, 16},
	{8, 3, 8, 10},
	{9, 3, 9, 6},
	{10, 3, 10, 4},
	{4, 4, 3, 40},
	{5, 5, 3, 30},
	{6, 6, 3, 20},
	{7, 7, 3, 16},
	{8, 8, 3, 10},
	{9, 9, 3, 6},
	{10, 10, 3, 4},
	{3, -1000, -1000, 2},
	{-1000, 3, -1000, 2},
	{-1000, -1000, 3, 2},
	{4, 3, 3, 40},
	{5, 3, 3, 30},
	{6, 3, 3, 20},
	{7, 3, 3, 16},
	{8, 3, 3, 10},
	{9, 3, 3, 6},
	{10, 3, 3, 4},
	{0, 3, 3, 2},
	{1, 3, 3, 2},
	{3, 4, 3, 40},
	{3, 5, 3, 30},
	{3, 6, 3, 20},
	{3, 7, 3, 16},
	{3, 8, 3, 10},
	{3, 9, 3, 6},
	{3, 10, 3, 4},
	{3, 0, 3, 2},
	{3, 1, 3, 2},
	{3, 3, 4, 40},
	{3, 3, 5, 30},
	{3, 3, 6, 20},
	{3, 3, 7, 16},
	{3, 3, 8, 10},
	{3, 3, 9, 6},
	{3, 3, 10, 4},
	{3, 3, 0, 2},
	{3, 3, 1, 2},
	{4, 4, 4, 20},
	{5, 5, 5, 15},
	{6, 6, 6, 10},
	{7, 7, 7, 8},
	{8, 8, 8, 5},
	{9, 9, 9, 3},
	{10, 10, 10, 2},
	{-1001, -1001, -1001, 1},
	{-1002, -1002, -1002, 1},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// NewPlate ...
func NewPlate(plateSize []int, scroll [][]int) ([]int, []int) {
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
func GameResult(plate []int) [][]int {
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
