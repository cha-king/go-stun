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

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}

	// stunConn, _ := multiplexConn(conn)

	stunConn, _ := NewVirtualConn(conn)

	client := stun.NewClient(stunConn)

	remoteAddr, err := client.BindRequest(addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", remoteAddr)
}
