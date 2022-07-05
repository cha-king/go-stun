package stun

import (
	"errors"
	"net"
)

type Client struct {
	conn net.PacketConn
}

func (c *Client) BindRequest(remoteAddr *net.UDPAddr) (net.UDPAddr, error) {
	m := newMessage(MessageClassRequest, MessageMethodBinding, nil)

	_, err := c.conn.WriteTo(m.encode(), remoteAddr)
	if err != nil {
		return net.UDPAddr{}, err
	}

	// TODO: Discard stun messages from unknown senders rather than error
	buf := make([]byte, 256)

	_, respAddr, err := c.conn.ReadFrom(buf)
	if err != nil {
		return net.UDPAddr{}, err
	}
	if remoteAddr.String() != respAddr.String() {
		return net.UDPAddr{}, errors.New("unknown sender")
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

	_, err := c.conn.WriteTo(m.encode(), remoteAddr)
	if err != nil {
		return err
	}
	return nil
}

func NewClient(conn net.PacketConn) *Client {
	return &Client{conn}
}
