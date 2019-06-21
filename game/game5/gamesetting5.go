package game5

// 貓下去

// Version ...
const Version = "1.0.3"

// GameIndex game sort id
const GameIndex = 5

// IsSingleLine game result only output one result
const isSingleLine = true

// BetRate ...
var betRate = []int64{1, 10, 100, 500, 1000, 2000, 5000, 10000, 20000, 500000}
var betRateLinkIndex = []int64{1, 2, 3, 4, 6}
var betRateDefaultIndex int64 = 1

// ScrollSize ...
var scrollSize = []int{1, 1, 1}

// NormalScroll ...
var normalScroll = [][]int{
    {8, 5, 0, 9, 8, 5, 6, 8, 5, 7, 0, 5, 6, 8, 5, 9, 8, 5, 10, 0, 10, 6, 5, 9, 8, 6, 9, 7, 5, 8, 10, 9, 7, 8, 9, 6, 5, 8, 9, 7, 9, 0, 9, 6, 7, 9, 8, 7, 6, 5, 8, 6, 7, 9, 8, 9, 7, 5, 8, 9, 5, 9, 8, 5, 7, 9, 8, 9, 0, 7, 5, 7, 10, 5, 8, 7, 9, 8, 7, 5, 10, 8, 9, 10, 7, 5, 9, 8, 6, 9, 5, 6},
    {8, 6, 9, 0, 5, 9, 7, 1, 8, 6, 9, 6, 7, 9, 5, 8, 10, 1, 6, 9, 10, 7, 5, 7, 8, 6, 5, 9, 8, 6, 7, 5, 8, 9, 6, 5, 8, 5, 6, 7, 9, 0, 7, 8, 7, 9, 8, 7, 1, 9, 8, 6, 7, 9, 8, 6, 1, 8, 6, 9, 8, 6, 8, 10, 7, 6, 8, 10, 9, 6, 7, 10, 7, 5, 1, 9, 10, 8, 9, 7, 8, 9, 10, 1, 10, 5, 10, 1, 7, 8, 5, 9},
    {6, 0, 8, 6, 9, 5, 8, 9, 5, 10, 0, 8, 5, 9, 7, 8, 9, 5, 8, 6, 8, 5, 9, 6, 5, 9, 8, 5, 10, 7, 9, 10, 5, 6, 9, 8, 7, 9, 7, 8, 6, 7, 8, 5, 7, 5, 10, 8, 7, 6, 9, 8, 5, 9, 7, 6, 7, 6, 0, 10, 8, 7, 10, 6, 7, 10, 0, 9, 5, 6, 7, 8, 9, 5, 10, 6, 10, 9, 7, 6, 9, 10, 7, 5, 10, 6, 5, 9, 10, 5, 8, 9},
}

// FreeGameTrigger free count equal FreeGameTrigger free game start
const FreeGameTrigger = 10

var freeGameWinRate = []int64{1, 2, 3, 4, 5}

// FreeScroll ...
var freeScroll = [][]int{
    {8, 7, 0, 10, 9, 5, 6, 8, 6, 10, 0, 9, 6, 8, 5, 9, 8, 9, 7, 9, 5, 6, 7, 9, 8, 6, 9, 7, 5, 8, 9, 8, 7, 8, 9, 6, 10, 9, 7, 9, 5, 9, 7, 6, 7, 9, 8, 0, 6, 5, 8, 6, 7, 9, 8, 9, 7, 9, 7, 9, 7, 9, 8, 9, 7, 8, 7, 9, 10, 0, 5, 7, 10, 5, 8, 7},
    {1},
    {8, 9, 8, 6, 9, 5, 6, 7, 9, 10, 0, 10, 6, 9, 5, 8, 9, 5, 8, 6, 10, 0, 10, 8, 7, 9, 8, 6, 7, 9, 10, 8, 6, 8, 9, 8, 7, 9, 7, 8, 6, 7, 8, 6, 9, 8, 7, 8, 7, 6, 9, 8, 7, 9, 7, 6, 7, 10, 6, 8, 9, 7, 9, 6, 10, 8, 10, 9, 5, 6, 7, 8, 9, 5, 10, 6},
}

// RespinScroll1 RTP 93
var respinScroll1 = []int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 4, 3, 2, 2, 2, 2, 2, 2, 2, 2}

// RespinScroll2 RTP 96
var respinScroll2 = []int{4, 3, 4, 2, 2, 2, 2, 2, 3, 2, 2, 3, 4, 4, 2, 3, 4, 3, 2, 2}

// RespinScroll3 RTP 99
var respinScroll3 = []int{3, 2, 4, 4, 3, 3, 4, 4, 2, 4, 4, 4, 2, 4, 4, 2, 4, 4, 3, 4}

// RespinSetting 1 RTP:93, 2 RTP:96, 3 RTP:99
var RespinSetting = 1

// Items item index
var items = []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Space Space item index
const space = 10
const wild1 = 0
const wild2 = 1

// ItemResults 0~10 item
// {item, item, item, result}
// -1000 any
// -1001 any 7
// -1002 any bar
// -100 bonus game 1
// -101 bonus game 2
var itemResults = [][]int{
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
