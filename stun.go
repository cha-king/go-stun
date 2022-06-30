package stun

import (
	"errors"
	"net"
)

func Bind(serverAddr *net.UDPAddr) (net.UDPAddr, error) {
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return net.UDPAddr{}, err
	}
	defer conn.Close()

	m := newMessage(MessageClassRequest, MessageMethodBinding, nil)

	_, err = conn.Write(m.encode())
	if err != nil {
		return net.UDPAddr{}, err
	}

	buf := make([]byte, 256)
	_, err = conn.Read(buf)
	if err != nil {
		return net.UDPAddr{}, err
	}

	r, err := decodeMessage(buf)
	if err != nil {
		return net.UDPAddr{}, err
	}

	if m.transactionId != r.transactionId {
		return net.UDPAddr{}, errors.New("transaction id mismatch")
	}

	attr := r.attributes[0]
	mAddr, ok := attr.(xorMappedAddress)
	if !ok {
		return net.UDPAddr{}, errors.New("unsupported attribute")
	}
	addr, err := mAddr.getAddress()
	if err != nil {
		return net.UDPAddr{}, err
	}

	return addr, nil
}
