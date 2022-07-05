package stun

import (
	"errors"
	"net"
)

type Client struct {
	conn net.Conn
}

func (c *Client) BindRequest() (net.UDPAddr, error) {
	m := newMessage(MessageClassRequest, MessageMethodBinding, nil)

	_, err := c.conn.Write(m.encode())
	if err != nil {
		return net.UDPAddr{}, err
	}

	buf := make([]byte, 256)
	_, err = c.conn.Read(buf)
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

func (c *Client) BindIndication(remoteAddr *net.UDPAddr) error {
	m := newMessage(MessageClassIndication, MessageMethodBinding, nil)

	_, err := c.conn.Write(m.encode())
	if err != nil {
		return err
	}
	return nil
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}
