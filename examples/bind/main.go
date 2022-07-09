package main

import (
	"fmt"
	"net"
	"time"

	"github.com/cha-king/go-stun"
)

const stunUrl = "stun.l.google.com:19302"

type fakeConn struct {
	conn net.PacketConn
}

func (c *fakeConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.conn.ReadFrom(p)
}
func (c *fakeConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.conn.WriteTo(p, addr)
}
func (c *fakeConn) Close() error {
	return c.conn.Close()
}
func (c *fakeConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *fakeConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *fakeConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *fakeConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func main() {
	addr, err := net.ResolveUDPAddr("udp4", stunUrl)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	mConn := &fakeConn{conn: conn}

	client := stun.NewClient(mConn)

	remoteAddr, err := client.BindRequest(addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", remoteAddr)
}
