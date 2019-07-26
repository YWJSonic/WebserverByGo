package gamerule

// 花開富貴

// Version ...
const Version = "0.0.1"

// GameIndex game sort id
const GameIndex = 6

// IsSingleLine game result only output one result
var isSingleLine = false

// IsAttachInit is init attach at player init.
// if init attach it will get all of game attach.
// not all or something  special setting set in "servicegame" file.
var IsAttachInit = false

// IsAttachSaveToDB is game attach save to db
var IsAttachSaveToDB = true

// WinScoreLimit game round win money limit
var WinScoreLimit int64

// WinBetRateLimit game round win rate limit
var WinBetRateLimit int64

// BetRate ...
var betRate = []int64{100, 500, 1000, 2000, 3000}
var betRateLinkIndex = []int64{0, 1, 2, 3, 4}
var betRateDefaultIndex int64 = 1

// ScrollSize ...
var scrollSize = []int{3, 3, 3, 3, 3}

// NormalScroll ...
var normalScroll = [][]int{
	{2, 7, 6, 9, 12, 5, 7, 8, 4, 10, 12, 6, 11, 8, 3, 7, 11, 1, 9, 12, 4, 7, 8, 5, 11, 10},
	{4, 12, 10, 0, 12, 11, 2, 9, 4, 7, 5, 12, 6, 11, 12, 1, 9, 11, 3, 7, 9, 2, 8, 1, 10, 6, 8, 5, 7},
	{1, 7, 2, 9, 4, 12, 6, 8, 1, 7, 4, 11, 5, 10, 1, 12, 5, 8, 11, 3, 9, 12, 2, 9, 11, 0, 8, 11},
	{8, 12, 10, 0, 12, 10, 3, 8, 4, 11, 2, 12, 6, 7, 1, 9, 4, 11, 3, 7, 8, 5, 7, 10, 1, 11, 6},
	{8, 2, 7, 4, 8, 5, 7, 9, 3, 10, 12, 6, 11, 8, 5, 10, 12, 1, 9, 7, 4, 11, 7, 12, 5, 11, 7, 6, 9, 7, 3, 10, 4, 11}}

// ScotterScroll RTP ...
var scotterScroll = [][]int{
	{2, 7, 6, 9, 12, 2, 7, 8, 4, 10, 12, 6, 11, 8, 3, 7, 11, 1, 9, 12, 4, 7, 8, 5, 11, 10},
	{4, 12, 10, 0, 12, 11, 2, 9, 4, 7, 5, 12, 6, 11, 12, 1, 9, 11, 3, 7, 9, 2, 8, 1, 10, 6, 8, 5, 7, 10},
	{1, 7, 2, 9, 4, 12, 6, 8, 1, 7, 4, 11, 5, 10, 1, 12, 5, 8, 11, 3, 9, 12, 2, 9, 11, 0, 9, 11, 10},
	{8, 12, 10, 0, 12, 10, 3, 8, 4, 11, 2, 12, 6, 7, 1, 9, 4, 11, 3, 7, 8, 5, 7, 10, 1, 11, 6, 9, 5},
	{8, 2, 7, 4, 8, 5, 7, 9, 3, 10, 12, 6, 11, 8, 5, 10, 12, 1, 9, 7, 4, 10, 7, 12, 5, 11, 7, 6, 9, 7, 3, 11, 4, 7}}

// ScotterGameSetting 1 RTP:93, 2 RTP:96, 3 RTP:99
var ScotterGameSetting = 1

// Items item    W1,F1,H1,H2,H3,H4,H5,L1,L2,L3, L4, L5, L6
var items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

// Space Space item index
const space = -1
const wild1 = 0
const scotter = 1
const scotterGameLimit = 3

// ItemResults 0~10 item
// {item, item, item, result}
// -1000 any
// -1001 any bar
// -100 bonus game 1
var itemResults = [][][]int{
	{},
	{},
	{},

	{{1, 1, 1, 5},
		{2, 2, 2, 50},
		{3, 3, 3, 35},
		{4, 4, 4, 30},
		{5, 5, 5, 20},
		{6, 6, 6, 15},
		{7, 7, 7, 10},
		{8, 8, 8, 10},
		{9, 9, 9, 10},
		{10, 10, 10, 10},
		{11, 11, 11, 5},
		{12, 12, 12, 5}},

	{{1, 1, 1, 1, 10},
		{2, 2, 2, 2, 100},
		{3, 3, 3, 3, 100},
		{4, 4, 4, 4, 100},
		{5, 5, 5, 5, 50},
		{6, 6, 6, 6, 35},
		{7, 7, 7, 7, 30},
		{8, 8, 8, 8, 20},
		{9, 9, 9, 9, 15},
		{10, 10, 10, 10, 15},
		{11, 11, 11, 11, 15},
		{12, 12, 12, 12, 10}},

	{{1, 1, 1, 1, 1, 50},
		{2, 2, 2, 2, 2, 1000},
		{3, 3, 3, 3, 3, 800},
		{4, 4, 4, 4, 4, 800},
		{5, 5, 5, 5, 5, 300},
		{6, 6, 6, 6, 6, 300},
		{7, 7, 7, 7, 7, 200},
		{8, 8, 8, 8, 8, 200},
		{9, 9, 9, 9, 9, 100},
		{10, 10, 10, 10, 10, 100},
		{11, 11, 11, 11, 11, 100},
		{12, 12, 12, 12, 12, 100}},
}

var normalWildWinRate = []int{1, 2, 3, 5, 8, 10, 15, 30, 40}
var normalWildWinRateWeightings = []int{1116, 60, 30, 15, 8, 5, 3, 2, 1}

var scotterGameSpinTime = []int{25, 20, 15, 13, 10, 6}
var scotterGameWildWinRate = [][]int64{
	{2, 3, 5},
	{3, 5, 8},
	{5, 8, 10},
	{8, 10, 15},
	{10, 15, 30},
	{15, 30, 40},
}
var scotterGameWildWinRateWeightings = [][]int{
	{16, 4, 2},
	{80, 34, 10},
	{23, 10, 4},
	{65, 5, 2},
	{80, 15, 4},
	{50, 15, 8},
}
var scotterH5SpecialWinRate = []int{50, 20, 15, 10, 5, 2}

var scotterH5SpecialWinRateWeightings = [][]int{
	{3, 5, 12, 14, 15, 15},
	{3, 7, 11, 12, 15, 21},
	{3, 7, 13, 13, 19, 21},
	{1, 2, 4, 5, 18, 26},
	{1, 2, 4, 6, 18, 21},
	{1, 5, 9, 13, 15, 20},
}

var scotterGameDefaultCombinationIndex = []int{0, 7, 14, 21, 28, 35}

// scotterGameMysteryIndexCombination WinRate, Spin
var scotterGameMysteryIndexCombination = [][]int{
	{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, // 0~5
	{0, 1}, {1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, // 6~11
	{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}, {5, 2}, // 12~17
	{0, 3}, {1, 3}, {2, 3}, {3, 3}, {4, 3}, {5, 3}, // 18~23
	{0, 4}, {1, 4}, {2, 4}, {3, 4}, {4, 4}, {5, 4}, // 24~29
	{0, 5}, {1, 5}, {2, 5}, {3, 5}, {4, 5}, {5, 5}, // 30~35
}
var scotterGameMysteryWeightings = []int{100, 76, 54, 46, 34, 20, 100, 100, 71, 60, 45, 30, 100, 100, 100, 85, 64, 40, 100, 100, 100, 100, 75, 44, 100, 100, 100, 100, 100, 53, 100, 100, 100, 100, 100, 100}
