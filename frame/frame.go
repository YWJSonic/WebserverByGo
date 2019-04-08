package frame

import "./code"

type Frame struct {
	Command code.Code   `json:"cmd"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
	Status  code.Code   `json:"status,omitempty"`
}

func New(command, status code.Code, data interface{}) Frame {
	return Frame{
		Command: command,
		Data:    data,
		Status:  status,
		Success: status == code.OK,
	}
}

func Prepare(v interface{}) Frame {
	return Frame{
		Data: v,
	}
}
