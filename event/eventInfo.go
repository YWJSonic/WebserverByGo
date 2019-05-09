package event

import "time"

// Event UpdateEvent struct
// 0: clock TimeSince >= DoTime Run
// 1: CrontabTime Run
type Event struct {
	Stat        int
	CrontabTime [5]int // month: -1~12 , dayOfmonth:-1~31, hour:-1~23, minute:-1~59, second:0~59. default -1 is everytime
	DoTime      time.Duration
	Do          func(interface{}) interface{}
}
