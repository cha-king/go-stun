package stun

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
)

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

	classBytes := (class&0b10)<<7 | (class&0b01)<<4

	methodBytes := (uint16(h.method) & 0x000F) | (uint16(h.method)&0x70)<<1 | (uint16(h.method)&0xF80)<<2

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

func DecodeHeader(headerBytes []byte) (header, error) {
	// Check first two bites
	if binary.BigEndian.Uint16(headerBytes[0:])>>14&0b11 != 0b00 {
		return header{}, errors.New("unable to parse header")
	}

	// Check magic cookie
	if cookie := binary.BigEndian.Uint32(headerBytes[4:]); cookie != MessageMagicCookie {
		return header{}, fmt.Errorf("unable to parse header: invalid magic cookie")
	}

	// TODO: Sanity check length
	length := binary.BigEndian.Uint16(headerBytes[2:])

	messageType := binary.BigEndian.Uint16(headerBytes[0:])
	class := uint8((messageType&0x010)>>4 | (messageType&0x100)>>7)
	method := uint8((messageType & 0x000F) | (messageType&0x00E0)>>1 | (messageType&0x1F00)>>2)

	var transactionId [transactionIdLength]byte
	copy(transactionId[:], headerBytes[8:20])

	return header{class, method, length, transactionId}, nil
}
