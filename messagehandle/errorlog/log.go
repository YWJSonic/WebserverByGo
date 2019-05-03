package errorlog

import (
	"fmt"

	"../../code"
)

// ErrorLogPrint ...
func ErrorLogPrint(msg string) {
	fmt.Print("Error " + msg)
}

// ErrorLogPrintf ...
func ErrorLogPrintf(msg string, a ...interface{}) {
	fmt.Printf("Error "+msg, a)
}

// ErrorLogPrintln ...
func ErrorLogPrintln(msg string, a ...interface{}) {
	fmt.Println("Error "+msg, a)
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
