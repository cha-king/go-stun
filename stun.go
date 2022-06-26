package main

import (
	"crypto/rand"
	"encoding/binary"
)

const (
	MessageClassRequest    uint8 = 0b00
	MessageClassIndication uint8 = 0b01
	MessageClassSuccess    uint8 = 0b10
	MessageClassError      uint8 = 0b11
)

const (
	MessageMethodBinding uint8 = 0x0001
)

const MessageMagicCookie uint32 = 0x2112A442

const transactionIdLength = 12

type header struct {
	class         uint8
	method        uint8
	length        uint16
	transactionId [transactionIdLength]byte
}

func (h header) encode() []byte {
	output := make([]byte, 20)
	class := uint16(h.class)

	classBytes := (class << 4 & 0x0F0) | (class << 8 & 0xF00)

	var methodBytes uint16
	methodBytes |= uint16(h.method) & 0x000F
	methodBytes |= uint16(h.method) << 1 & 0x00F0
	methodBytes |= uint16(h.method) << 2 & 0xFF00

	typeBytes := classBytes | methodBytes

	binary.BigEndian.PutUint16(output[0:2], typeBytes)
	binary.BigEndian.PutUint16(output[2:4], h.length)
	binary.BigEndian.PutUint32(output[4:8], MessageMagicCookie)
	binary.BigEndian.PutUint32(output[8:12], binary.BigEndian.Uint32(h.transactionId[0:4]))
	binary.BigEndian.PutUint32(output[12:16], binary.BigEndian.Uint32(h.transactionId[4:8]))
	binary.BigEndian.PutUint32(output[16:20], binary.BigEndian.Uint32(h.transactionId[8:12]))

	return output
}

func newHeader(class uint8, method uint8) header {
	var transactionId [transactionIdLength]byte
	randBytes := make([]byte, transactionIdLength)
	_, err := rand.Read(randBytes)
	if err != nil {
		panic(err)
	}
	copy(transactionId[:], randBytes)
	return header{class, method, 0, transactionId}
}
