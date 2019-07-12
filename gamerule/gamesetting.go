package gamerule

//

// Version ...
const Version = "0.0.1"

// GameIndex game sort id
const GameIndex = 8

// IsSingleLine game result only output one result
var isSingleLine = false

// IsAttachSaveToDB is game attach save to db
var IsAttachSaveToDB = true

// WinScoreLimit game round win money limit
var WinScoreLimit int64

// WinBetRateLimit game round win rate limit
var WinBetRateLimit int64

// BetRate ...
var betRate = []int64{1000, 5000, 10000, 20000, 30000}
var betRateLinkIndex = []int64{0, 1, 2, 3, 4}
var betRateDefaultIndex int64 = 1

// ScrollSize ...
var scrollSize = []int{3, 3, 3, 3, 3}

// NormalScroll ...
var normalScroll = [][]int{
	{0, 0, 2, 3, 0, 5, 6, 0, 8},
	{0, 1, 0, 5, 4, 5, 6, 7, 8},
	{0, 1, 5, 3, 4, 0, 6, 5, 8},
	{0, 7, 2, 0, 4, 5, 6, 7, 8},
	{0, 1, 2, 0, 4, 7, 0, 7, 8},
}

// RespinScroll1 RTP ...
var respinScroll1 = []int{4, 8, 4, 8, 4, 8, 1, 8, 5, 7, 6, 8, 7, 7, 8, 7, 8, 7, 8, 7, 8, 5, 8, 7, 6, 7, 8, 5, 7, 7, 7, 8, 6, 7, 8, 7, 7, 6, 7, 5, 8, 2, 8, 7, 8, 6, 5, 8, 7, 7, 7, 6, 8, 6, 7, 5, 6, 8, 3, 8, 3, 8, 3, 8, 6, 7, 6, 7, 7, 1, 8, 1, 8, 1, 8, 4, 8, 5, 7, 7, 6, 8, 6, 7, 5, 8, 6, 7, 7, 6, 8, 6, 7, 8, 5, 8, 6, 4, 5, 7, 8, 5, 4, 8, 4, 4, 7, 6, 7, 8, 5, 6, 6, 5, 8, 2, 8, 6, 8, 6, 5, 6, 7, 5, 5, 6, 7, 6, 8, 4, 8, 6, 7, 7, 6, 5, 6, 6, 5, 7, 6, 7, 7, 7, 6, 8, 6, 7, 7, 6, 7, 8, 7, 7, 5, 6, 7, 7, 5, 6}

// RespinSetting 1 RTP:93, 2 RTP:96, 3 RTP:99
var RespinSetting = 1

// index 0 x5, 1:x3, 2:x2
var jackPortTex = []float64{0.007, 0.006, 0.004}
var jackPartWinRate = []int{30, 45, 75}

// Items item index
var items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8}

// Space Space item index
const space = -1
const wild1 = 0

// 2x
// const wild2 = 1

// // 3x
// const wild3 = 2

// // 5x
// const wild4 = 3

// ResultRateArray for client count money value
var resultRateArray = []int{
	75, 45, 30, 10, 5, 3, 2, 1,
}

// ItemResults 0~10 item
// {item, item, item, result}
// -1000 any
// -1001 any bar
// -100 bonus game 1
var itemResults = [][][]int{
	{},
	{},
	{},

	{{0, 0, 0, 97},
		{-1000, 1, 1, -100},
		{2, 2, 2, 77},
		{3, 3, 3, 67},
		{4, 4, 4, 57},
		{5, 5, 5, 47},
		{6, 6, 6, 37},
		{7, 7, 7, 27},
		{8, 8, 8, 17}},

	{{0, 0, 0, 0, -100},
		{1, 1, 1, 1, 88},
		{-1000, -1000, -1000, 2, 78},
		{3, 3, 3, 3, 68},
		{4, 4, 4, 4, 58},
		{5, 5, 5, 5, 48},
		{6, 6, 6, 6, 38},
		{7, 7, 7, 7, 28},
		{8, 8, 8, 8, 18}},

	{{1, 1, 1, 1, 1, -101},
		{-1001, -1001, -1001, -1001, -1001, 79},
		{3, 3, 3, 3, 3, 69},
		{4, 4, 4, 4, 4, 59},
		{5, 5, 5, 5, 5, 49},
		{6, 6, 6, 6, 6, 39},
		{7, 7, 7, 7, 7, 29},
		{8, 8, 8, 8, 8, 19}},
}

var lineMap = [][]int{
	{1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0},
	{2, 2, 2, 2, 2},
	{0, 1, 2, 1, 0},
	{2, 1, 0, 1, 2},
	{1, 0, 0, 0, 1},
	{1, 2, 2, 2, 1},
	{0, 0, 1, 2, 2},
	{2, 2, 1, 0, 0},
	{1, 2, 1, 0, 1},
	{1, 0, 1, 2, 1},
	{0, 1, 1, 1, 0},
	{2, 1, 1, 1, 2},
	{0, 1, 0, 1, 0},
	{2, 1, 2, 1, 2},
	{1, 1, 0, 1, 1},
	{1, 1, 2, 1, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 2, 2, 2, 0},
}

var jackPortResults = [][]int{
	{0, 1, 0, -101},
	{0, 2, 0, -102},
	{0, 3, 0, -103},
}

var respinitemResults = [][]int{
	{0, 4, 0, 10},
	{0, 5, 0, 5},
	{0, 6, 0, 3},
	{0, 7, 0, 2},
}

var symbolGroup = map[int][]int{
	-1000: []int{0, 1, 2, 3, 4, 5, 6, 7},
	-1001: []int{5, 6, 7},
}
var spWhildWinRate = []int64{2, 3, 5}
