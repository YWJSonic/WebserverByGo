package gamelogic

// GetGameResult All game result interface
func GetGameResult(gameid string, bet int) interface{} {

	switch gameid {
	case "slot":
		gameOutput([]int{0, 1, 2, 3}, []int{3, 3, 3, 3, 3})
	default:
	}
	return "str"

}
