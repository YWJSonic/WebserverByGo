package errorlog

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGo/code"
)

// IsAddTimeFlag log title add time flag
var IsAddTimeFlag = true

// log Tag
const (
	Log      = "Log"
	ErrorLog = "Error"
)

// LogPrint ...
func LogPrint(msg string) {
	if IsAddTimeFlag {
		fmt.Printf("%s %s %s", time.Now().Format(time.Stamp), Log, msg)
	} else {
		fmt.Printf("%s %s", Log, msg)
	}
}

// LogPrintf ...
func LogPrintf(msg string, a ...interface{}) {
	var str string
	if IsAddTimeFlag {
		str = fmt.Sprintf("%s %s %s", time.Now().Format(time.Stamp), Log, msg)
	} else {
		str = fmt.Sprintf("%s %s", Log, msg)
	}
	fmt.Printf(str, a...)
}

// LogPrintln ...
func LogPrintln(msg string, a ...interface{}) {
	var tmp []interface{}
	if IsAddTimeFlag {
		tmp = append(tmp, time.Now().Format(time.Stamp))
	}
	tmp = append(tmp, Log)
	tmp = append(tmp, msg)
	tmp = append(tmp, a...)
	fmt.Println(tmp...)
}

// ErrorLogPrint ...
func ErrorLogPrint(msg string) {
	if IsAddTimeFlag {
		fmt.Printf("%s %s %s", time.Now().Format(time.Stamp), ErrorLog, msg)
	} else {
		fmt.Printf("%s %s", ErrorLog, msg)
	}
}

// ErrorLogPrintf ...
func ErrorLogPrintf(msg string, a ...interface{}) {
	var str string
	if IsAddTimeFlag {
		str = fmt.Sprintf("%s %s %s", time.Now().Format(time.Stamp), ErrorLog, msg)
	} else {
		str = fmt.Sprintf("%s %s", ErrorLog, msg)
	}
	fmt.Printf(str, a...)
}

// ErrorLogPrintln ...
func ErrorLogPrintln(msg string, a ...interface{}) {
	var tmp []interface{}
	if IsAddTimeFlag {
		tmp = append(tmp, time.Now().Format(time.Stamp))
	}
	tmp = append(tmp, ErrorLog)
	tmp = append(tmp, msg)
	tmp = append(tmp, a...)
	fmt.Println(tmp...)
}

// ErrorMsg ...
type ErrorMsg struct {
	ErrorCode code.Code
	// MsgNum    int8
	Msg string
}

// New default Error Message
func New() ErrorMsg {
	return ErrorMsg{
		ErrorCode: code.OK,
		// MsgNum:    0,
		Msg: "",
	}
}
