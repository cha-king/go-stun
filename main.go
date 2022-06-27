package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func runClient() {
	h := newHeader(MessageClassRequest, MessageMethodBinding)

	addr, err := net.ResolveUDPAddr("udp4", "stun.l.google.com:19302")
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	headerBytes := h.encode()
	_, err = conn.Write(headerBytes)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}

	_, err = decode(buf)
	if err != nil {
		panic(err)
	}

	length := binary.BigEndian.Uint16(buf[2:])

	attributes := buf[20 : 20+length]
	attributeType := binary.BigEndian.Uint16(attributes[0:])
	attributeLength := binary.BigEndian.Uint16(attributes[2:])

	if attributeType != 0x0020 {
		return
	}

	xorMappedAddress := attributes[4 : 4+attributeLength]
	family := binary.BigEndian.Uint16(xorMappedAddress[0:])
	if family != 0x01 {
		return
	}
	xPort := binary.BigEndian.Uint16(xorMappedAddress[2:])
	port := xPort ^ uint16(MessageMagicCookie>>16)
	fmt.Printf("Port: %d\n", port)

	xAddress := binary.BigEndian.Uint32(xorMappedAddress[4:])
	addressBytes := xAddress ^ MessageMagicCookie
	addressSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(addressSlice, addressBytes)
	ip := net.IP(addressSlice)
	fmt.Printf("Address: %s\n", ip)
}

func runServer() {
	addr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:8000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening")
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buf := make([]byte, 548)
	_, remoteAddr, err := conn.ReadFromUDP(buf)
	if err != nil {
		panic(err)
	}

	// Check first two bits
	if binary.BigEndian.Uint16(buf[0:])>>14&0b11 != 0b00 {
		return
	}

	// Check magic cookie
	if binary.BigEndian.Uint32(buf[4:]) != MessageMagicCookie {
		return
	}
	fmt.Println("Valid!")
	_, err = conn.WriteToUDP([]byte("Hello world"), remoteAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent")
}

func main() {
	arg := os.Args[1]
	if arg == "client" {
		runClient()
	} else if arg == "server" {
		runServer()
	}
}
