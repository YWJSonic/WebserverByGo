package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	time.Microsecond
	address := net.TCPAddr{
		IP:   net.ParseIP("192.168.0.105"),
		Port: 8000,
		Zone: "",
	}
	listener, err := net.ListenTCP("tcp", &address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Other Address: ", conn.RemoteAddr())
	}
}
