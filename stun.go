package stun

import (
	"errors"
	"net"
)

type Client struct {
	conn net.PacketConn
}

func (c *Client) BindRequest(addr net.Addr) (net.UDPAddr, error) {
	m := newMessage(MessageClassRequest, MessageMethodBinding, nil)

	_, err := c.conn.WriteTo(m.encode(), addr)
	if err != nil {
		return net.UDPAddr{}, err
	}

	// TODO: Handle receiving from wrong IP
	buf := make([]byte, 256)
	_, _, err = c.conn.ReadFrom(buf)
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
	localaddr, err := mAddr.getAddress()
	if err != nil {
		return net.UDPAddr{}, err
	}

	return localaddr, nil
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
