package main

import (
	"fmt"
	"net"

	"github.com/cha-king/go-stun"
)

const (
	stunUrl  = "stun.l.google.com:19302"
	localUrl = "0.0.0.0:8000"
	peerUrl  = "35.87.220.143:50000"
)

func wrapConnection(conn net.Conn) net.Conn {
	leftConn, rightConn := net.Pipe()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := rightConn.Read(buf)
			if err != nil {
				panic(err)
			}
			conn.Write(buf[:n])
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				panic(err)
			}
			rightConn.Write(buf[:n])
		}
	}()

	return leftConn
}

func main() {
	localAddr, err := net.ResolveUDPAddr("udp4", localUrl)
	if err != nil {
		panic(err)
	}

	stunAddr, err := net.ResolveUDPAddr("udp4", stunUrl)
	if err != nil {
		panic(err)
	}

	peerAddr, err := net.ResolveUDPAddr("udp4", peerUrl)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", localAddr, stunAddr)
	if err != nil {
		panic(err)
	}

	wrappedConn := wrapConnection(conn)

	client := stun.NewClient(wrappedConn)

	resp, err := client.BindRequest()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", resp)

	err = client.BindIndication(peerAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening..")

	buf := make([]byte, 1024)
	_, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}

	_, err = conn.WriteToUDP(buf, addr)
	if err != nil {
		panic(err)
	}
}
