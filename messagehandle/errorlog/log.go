package errorlog

import (
	"fmt"

	"../../code"
)

// LogPrint ...
func LogPrint(msg string) {
	fmt.Print("Log " + msg)
}

// LogPrintf ...
func LogPrintf(msg string, a ...interface{}) {
	fmt.Printf("Log "+msg, a...)
}

// LogPrintln ...
func LogPrintln(msg string, a ...interface{}) {
	var tmp []interface{}
	tmp = append(tmp, "Log")
	tmp = append(tmp, msg)
	tmp = append(tmp, a...)
	fmt.Println(tmp...)
}

// ErrorLogPrint ...
func ErrorLogPrint(msg string) {
	fmt.Print("Error " + msg)
}

// ErrorLogPrintf ...
func ErrorLogPrintf(msg string, a ...interface{}) {
	fmt.Printf("Error "+msg, a...)
}

// ErrorLogPrintln ...
func ErrorLogPrintln(msg string, a ...interface{}) {
	var tmp []interface{}
	tmp = append(tmp, "Error")
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
