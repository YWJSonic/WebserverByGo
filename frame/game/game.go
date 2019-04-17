package game

// GameID ...
var GameID map[string]int

func init() {

}
func gameTypeInit() {
	GameID = make(map[string]int)
	GameID["Slot"] = 1
	// GameTypeID["Slot2"] = 2
	// sicbo
	// GameTypeID["sicbo"] = 1001
}
