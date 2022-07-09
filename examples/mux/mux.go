package main

import (
	"net"
	"time"
)

type virtualConn struct {
	conn net.PacketConn
}

func (c *virtualConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	panic("method not implemented")
}
func (c *virtualConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	panic("method not implemented")
}
func (c *virtualConn) Close() error {
	panic("method not implemented")
}
func (c *virtualConn) LocalAddr() net.Addr {
	panic("method not implemented")
}
func (c *virtualConn) SetDeadline(t time.Time) error {
	panic("method not implemented")
}
func (c *virtualConn) SetReadDeadline(t time.Time) error {
	panic("method not implemented")
}
func (c *virtualConn) SetWriteDeadline(t time.Time) error {
	panic("method not implemented")
}

func multiplexConn(conn net.PacketConn) (stunConn net.PacketConn, appConn net.PacketConn) {
	return &virtualConn{conn}, &virtualConn{conn}
}
