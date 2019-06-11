package gamesystem

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// ResultMap game result base info
func ResultMap(scrollIndex, plate []int, scores int64, islink bool) map[string]interface{} {
	result := make(map[string]interface{})

	result["plateindex"] = scrollIndex
	result["plate"] = plate
	result["scores"] = scores
	if islink {
		result["islink"] = 1
	} else {
		result["islink"] = 0
	}
	return result
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
