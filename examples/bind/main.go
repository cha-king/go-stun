package main

import (
	"fmt"
	"net"

	"github.com/cha-king/go-stun"
)

const stunUrl = "stun.l.google.com:19302"

func main() {
	addr, err := net.ResolveUDPAddr("udp4", stunUrl)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := stun.NewClient(conn)

	remoteAddr, err := client.BindRequest()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", remoteAddr)
}
