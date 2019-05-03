package gamelogic

import (
	"fmt"
	"math/rand"
	"time"
)

// GameOutput all game result output func
func gameOutput(items []int, wheelsSize []int) interface{} {

	result := newResult(items, wheelsSize)

	for _, value := range result {
		fmt.Println(value)
	}

	logicdata := Slot243Logic(items, result, GameLogicOption{
		While: 0,
		Sttel: -1,
	})

	for key, value := range logicdata {
		fmt.Println(key, value)
	}
	fmt.Println("---------------")
	return converTosingleLine(result)
}

func converTosingleLine(result [][]int) []int {
	var tmp []int
	for _, value := range result {
		tmp = append(tmp, value...)
	}
	return tmp
}
func newResult(items []int, wheelsSize []int) [][]int {
	var resultRow [][]int
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for _, wheelsize := range wheelsSize {
		// resultCol := make([]int, wheelsize)
		var resultCol []int
		for weelindex := 0; weelindex < wheelsize; weelindex++ {

			resultCol = append(resultCol, r1.Intn(len(items)))
		}
		resultRow = append(resultRow, resultCol)
	}
	// resultRow = [][]int{{0, 2, 2},
	// 	{1, 2, 1},
	// 	{1, 2, 3},
	// 	{1, 2, 0},
	// 	{2, 0, 0}}
	return resultRow
}
func rateResult(items []int, result map[int][][]int) {
	// rateInfo := map[int][]slotRate{
	// 	0: []slotRate{slotRate{0, 2, 200, 0, 0}, slotRate{0, 3, 300, 0, 0}, slotRate{0, 4, 400, 0, 0}},
	// 	1: []slotRate{slotRate{1, 3, 30, 0, 0}, slotRate{1, 4, 40, 0, 0}, slotRate{1, 5, 50, 0, 0}},
	// 	2: []slotRate{slotRate{2, 3, 3, 0, 0}, slotRate{2, 4, 4, 0, 0}, slotRate{2, 5, 5, 0, 0}},
	// 	3: []slotRate{slotRate{3, 2, 1, 0, 0}},
	// }

}

// Slot243Logic ...
func Slot243Logic(items []int, normalGame [][]int, option GameLogicOption) map[int][][]int {

	itemCount := make(map[int][][]int)
	ColIndex := 0
	AllIndex := 0

	//		// col
	//	row // 0 3 6 9  12
	//		// 1 4 7 10 13
	//		// 2 5 8 11 14
	for _, ColItem := range normalGame {
		RowIndex := 0
		for _, RowItem := range ColItem {
			if RowItem == option.While {
				if ColIndex == 0 {
					for _, itemIndex := range items {
						if len(itemCount[RowItem])+len(itemCount[option.While]) >= ColIndex {
							// if !contains(itemCount, itemIndex) {
							// 	itemCount[itemIndex] = [][]int{}
							// }
							// if len(itemCount[itemIndex]) >= ColIndex {
							// 	itemCount[itemIndex] = append(itemCount[itemIndex], []int{})
							// }

							if ColIndex-len(itemCount[itemIndex]) > 0 {
								continue
							}
							formatCheck(&itemCount, itemIndex, ColIndex)
							itemCount[itemIndex][ColIndex] = append(itemCount[itemIndex][ColIndex], AllIndex)

						}
					}
				} else {
					for itemIndex := range itemCount {
						// if len(itemCount[itemIndex]) >= ColIndex {
						// 	itemCount[itemIndex] = append(itemCount[itemIndex], []int{})
						// }
						if ColIndex-len(itemCount[itemIndex]) > 0 {
							continue
						}
						formatCheck(&itemCount, itemIndex, ColIndex)
						if len(itemCount[itemIndex])+len(itemCount[option.While]) >= ColIndex {
							itemCount[itemIndex][ColIndex] = append(itemCount[itemIndex][ColIndex], AllIndex)
						}
					}
				}
			} else {
				if len(itemCount[RowItem])+len(itemCount[option.While]) >= ColIndex {
					// if ColIndex == 0 && !contains(itemCount, RowItem) {
					// 	itemCount[RowItem] = [][]int{}
					// }
					// if len(itemCount[RowItem]) > ColIndex {
					// 	itemCount[RowItem] = append(itemCount[RowItem], []int{})
					// }

					if ColIndex-len(itemCount[RowItem]) > 0 {
						continue
					}
					formatCheck(&itemCount, RowItem, ColIndex)
					itemCount[RowItem][ColIndex] = append(itemCount[RowItem][ColIndex], AllIndex)
				}
			}
			RowIndex++
			AllIndex++
		}
		ColIndex++
	}
	return itemCount
}
func formatCheck(itemCount *map[int][][]int, RowItem int, ColIndex int) {
	if ColIndex == 0 && !contains(*itemCount, RowItem) {
		(*itemCount)[RowItem] = [][]int{}
	}
	if !(len((*itemCount)[RowItem]) > ColIndex) {
		(*itemCount)[RowItem] = append((*itemCount)[RowItem], []int{})
	}
}

func contains(array map[int][][]int, target interface{}) bool {

	if _, ok := array[target.(int)]; ok {
		return true
	}
	return false
}

type slotRate struct {
	ItemID       int
	Count        int
	ResultMmoney int
	BonusGame    int
	FreeGame     int
}

// GameLogicOption ...
type GameLogicOption struct {
	While int
	Sttel int
}
