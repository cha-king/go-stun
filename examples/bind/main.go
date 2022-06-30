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

	remoteAddr, err := stun.Bind(addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", remoteAddr)
}
