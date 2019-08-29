package gamerule

//

// Version ...
const Version = "0.0.1"

// GameIndex game sort id
const GameIndex = 8

// IsAttachInit ...
const IsAttachInit = false

// IsSingleLine game result only output one result
var isSingleLine = true

// IsAttachSaveToDB is game attach save to db
var IsAttachSaveToDB = true

// WinScoreLimit game round win money limit
var WinScoreLimit int64

// WinBetRateLimit game round win rate limit
var WinBetRateLimit int64

// BetRate ...
var betRate []int64
var betRateLinkIndex []int64
var betRateDefaultIndex int64

// ScrollSize ...
var scrollSize []int

// NormalScroll ...
var normalScroll [][]int

// FreeScroll ...
var freeScroll [][]int

// RespinScroll1 RTP ...
var respinScroll1 []int

// RespinSetting 1 RTP:93, 2 RTP:96, 3 RTP:99
var RespinSetting = 1

// Items item index
var items []int

// Space Space item index
const space = 8
const wild1 = 0
const scotter1 = 1

// ResultRateArray for client count money value
var resultRateArray []int

// ItemResults 0~10 item
// {item, item, item, result}
// -1000 any
// -1001 any bar
// -100 bonus game 1
var itemResults [][]int
var jackPortResults [][]int
var respinitemResults [][]int
var symbolGroup map[int][]int
var spWhildWinRate []int64
