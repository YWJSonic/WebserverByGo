package game5

import (
	"math/rand"
	"time"
)

// 貓下去

// GameIndex game sort id
const GameIndex = 5

// IsSingleLine game result only output one result
const IsSingleLine = true

// ScrollSize ...
var ScrollSize = []int{1, 1, 1}

// NormalScroll ...
var NormalScroll = [][]int{
	{8, 5, 0, 9, 8, 5, 6, 8, 5, 7, 0, 5, 6, 8, 5, 9, 8, 5, 10, 0, 10, 6, 5, 9, 8, 6, 9, 7, 5, 8, 10, 9, 7, 8, 9, 6, 5, 8, 9, 7, 9, 0, 9, 6, 7, 9, 8, 7, 6, 5, 8, 6, 7, 9, 8, 9, 7, 5, 8, 9, 5, 9, 8, 5, 7, 9, 8, 9, 0, 7, 5, 7, 10, 5, 8, 7, 9, 8, 7, 5, 10, 8, 9, 10, 7, 5, 9, 8, 6, 9, 5, 6},
	{8, 6, 9, 0, 5, 9, 7, 1, 8, 6, 9, 6, 7, 9, 5, 8, 10, 1, 6, 9, 10, 7, 5, 7, 8, 6, 5, 9, 8, 6, 7, 5, 8, 9, 6, 5, 8, 5, 6, 7, 9, 0, 7, 8, 7, 9, 8, 7, 1, 9, 8, 6, 7, 9, 8, 6, 1, 8, 6, 9, 8, 6, 8, 10, 7, 6, 8, 10, 9, 6, 7, 10, 7, 5, 1, 9, 10, 8, 9, 7, 8, 9, 10, 1, 10, 5, 10, 1, 7, 8, 5, 9},
	{6, 0, 8, 6, 9, 5, 8, 9, 5, 10, 0, 8, 5, 9, 7, 8, 9, 5, 8, 6, 8, 5, 9, 6, 5, 9, 8, 5, 10, 7, 9, 10, 5, 6, 9, 8, 7, 9, 7, 8, 6, 7, 8, 5, 7, 5, 10, 8, 7, 6, 9, 8, 5, 9, 7, 6, 7, 6, 0, 10, 8, 7, 10, 6, 7, 10, 0, 9, 5, 6, 7, 8, 9, 5, 10, 6, 10, 9, 7, 6, 9, 10, 7, 5, 10, 6, 5, 9, 10, 5, 8, 9},
}

// FreeGameTrigger free count equal FreeGameTrigger free game start
const FreeGameTrigger = 6

// FreeScroll ...
var FreeScroll = [][]int{
	{8, 7, 0, 10, 9, 5, 6, 8, 6, 10, 0, 9, 6, 8, 5, 9, 8, 9, 7, 9, 5, 6, 7, 9, 8, 6, 9, 7, 5, 8, 9, 8, 7, 8, 9, 6, 10, 9, 7, 9, 5, 9, 7, 6, 7, 9, 8, 0, 6, 5, 8, 6, 7, 9, 8, 9, 7, 9, 7, 9, 7, 9, 8, 9, 7, 8, 7, 9, 10, 0, 5, 7, 10, 5, 8, 7},
	{1},
	{8, 9, 8, 6, 9, 5, 6, 7, 9, 10, 0, 10, 6, 9, 5, 8, 9, 5, 8, 6, 10, 0, 10, 8, 7, 9, 8, 6, 7, 9, 10, 8, 6, 8, 9, 8, 7, 9, 7, 8, 6, 7, 8, 6, 9, 8, 7, 8, 7, 6, 9, 8, 7, 9, 7, 6, 7, 10, 6, 8, 9, 7, 9, 6, 10, 8, 10, 9, 5, 6, 7, 8, 9, 5, 10, 6},
}

// RespinScroll1 RTP 93
var RespinScroll1 = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 3, 2, 2, 2, 2, 2, 2, 2, 2}

// RespinScroll2 RTP 96
var RespinScroll2 = []int{4, 3, 4, 2, 2, 2, 2, 2, 3, 2, 2, 3, 4, 4, 2, 3, 4, 3, 2, 2}

// RespinScroll3 RTP 99
var RespinScroll3 = []int{3, 2, 4, 4, 3, 3, 4, 4, 2, 4, 4, 4, 2, 4, 4, 2, 4, 4, 3, 4}

// RespinSetting 1 RTP:93, 2 RTP:96, 3 RTP:99
const RespinSetting = 1

// Items item index
var Items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Space Space item index
const Space = 10

// ItemResults 0~10 item
// {item, item, item, result}
// -1000 any
// -1001 any 7
// -1002 any bar
// -100 bonus game 1
// -101 bonus game 2
var ItemResults = [][]int{
	{-1000, 1, 0, -100},
	{0, 0, 0, 20},
	{5, 5, 5, 10},
	{6, 6, 6, 6},
	{7, 7, 7, 4},
	{8, 8, 8, 3},
	{9, 9, 9, 2},
	{-1001, -1001, -1001, 2},
	{-1002, -1002, -1002, 1},
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RespuinScroll ...
func RespuinScroll() []int {
	if RespinSetting == 1 {
		return RespinScroll1
	} else if RespinSetting == 2 {
		return RespinScroll2
	} else {
		return RespinScroll3
	}
}

// NoprmalPlate other data 0: is resprin, 1: is count freesprint
func NoprmalPlate() ([]int, []int, []int) {
	var otherdata []int

	ScrollIndex, plate := NewPlate([]int{1, 1, 1}, NormalScroll)
	if plate[0] != 10 && plate[1] == 1 && plate[2] == 0 {
		otherdata = append(otherdata, 1)
	} else {
		otherdata = append(otherdata, 0)
	}
	if plate[1] == 1 {
		otherdata = append(otherdata, 1)
	} else {
		otherdata = append(otherdata, 0)
	}

	return ScrollIndex, plate, otherdata
}

// FreePlate ...
func FreePlate() ([]int, []int, []int) {
	var otherdata []int
	ScrollIndex, plate := NewPlate([]int{1, 1, 1}, FreeScroll)
	return ScrollIndex, plate, otherdata
}

// RespinPlate ...
func RespinPlate() ([]int, []int, []int) {
	var otherdata []int
	ScrollIndex, plate := NewPlate([]int{1}, [][]int{RespuinScroll()})
	return ScrollIndex, plate, otherdata
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
			if IsSingleLine {
				return result
			}
		}
	}
	return result
}

// RespinResult ...
func RespinResult(plate []int) [][]int {
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
	if Item == 5 || Item == 6 {
		return true
	}
	return false
}

func isAnyBay(Item int) bool {
	if Item == 6 || Item == 7 || Item == 8 || Item == 9 {
		return true
	}
	return false
}
