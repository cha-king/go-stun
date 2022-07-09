package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cha-king/go-stun"
)

const stunUrl = "stun.l.google.com:19302"
const remoteUrl = "35.87.255.115:8000"

func main() {
	addr, err := net.ResolveUDPAddr("udp4", stunUrl)
	if err != nil {
		panic(err)
	}

	remoteAddr, err := net.ResolveUDPAddr("udp4", remoteUrl)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}

	// stunConn, _ := multiplexConn(conn)

	stunConn, appConn := NewVirtualConn(conn)

	client := stun.NewClient(stunConn)

	go func() {
		for {
			remoteAddr, err := client.BindRequest(addr)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%+v\n", remoteAddr)
			time.Sleep(1 * time.Second)
		}
	}()

	// Normally, remote addr would be determined via ICE
	// and this write would be done as another bind request to STUN peer
	appConn.WriteTo([]byte("foo"), remoteAddr)

	buf := make([]byte, 1024)
	for {
		n, addr, err := appConn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Echoing %s\n", buf[:n])
		appConn.WriteTo(buf[:n], addr)
	}
}
