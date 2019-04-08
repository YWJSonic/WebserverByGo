package src

import "../frame"

type Client interface {
	Close()

	Reciver() <-chan frame.Frame
}
