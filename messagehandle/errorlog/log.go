package errorlog

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGoGame5/code"
)

var (
	// IsAddTimeFlag log title add time flag
	IsAddTimeFlag = true
	// IsPrintTipLog Tip log
	IsPrintTipLog = true
	// IsPrintLog Debug log
	IsPrintLog = true
	// IsPrintErrorLog Error log
	IsPrintErrorLog = true
)

// log Tag
const (
	Log      = "Log"
	ErrorLog = "Error"
)

// TipLogPirnt ...
func TipLogPirnt(msg string) {
	print(Log, msg)
}

// TipPrintf ...
func TipPrintf(msg string, a ...interface{}) {
	printf(Log, msg, a...)

}

// TipPrintln ...
func TipPrintln(msg string, a ...interface{}) {
	println(Log, msg, a...)

}

// LogPrint ...
func LogPrint(msg string) {

	if !IsPrintLog {
		return
	}

	print(Log, msg)
}

// LogPrintf ...
func LogPrintf(msg string, a ...interface{}) {

	if !IsPrintLog {
		return
	}

	printf(Log, msg, a...)
}

// LogPrintln ...
func LogPrintln(msg string, a ...interface{}) {

	if !IsPrintLog {
		return
	}
	println(Log, msg, a...)
}

// ErrorLogPrint ...
func ErrorLogPrint(msg string) {
	print(ErrorLog, msg)
}

// ErrorLogPrintf ...
func ErrorLogPrintf(msg string, a ...interface{}) {
	printf(ErrorLog, msg, a...)
}

// ErrorLogPrintln ...
func ErrorLogPrintln(msg string, a ...interface{}) {
	println(ErrorLog, msg, a...)
}

func print(logtype, msg string, a ...interface{}) {
	if IsAddTimeFlag {
		fmt.Printf("%s %s %s", time.Now().Format(time.Stamp), logtype, msg)
	} else {
		fmt.Printf("%s %s", logtype, msg)
	}
}
func printf(logtype, msg string, a ...interface{}) {
	if IsAddTimeFlag {
		msg = fmt.Sprintf("%s %s %s", time.Now().Format(time.Stamp), logtype, msg)
	} else {
		msg = fmt.Sprintf("%s %s", logtype, msg)
	}
	fmt.Printf(msg, a...)
}
func println(logtype, msg string, a ...interface{}) {
	var tmp []interface{}
	if IsAddTimeFlag {
		tmp = append(tmp, time.Now().Format(time.Stamp))
	}
	tmp = append(tmp, logtype)
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
