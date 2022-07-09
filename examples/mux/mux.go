package main

import (
	"net"
	"time"
)

type virtualConn struct {
	conn net.PacketConn
}

func (c *virtualConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.conn.ReadFrom(p)
}

func (c *virtualConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.conn.WriteTo(p, addr)
}
func (c *virtualConn) Close() error {
	return c.conn.Close()
}
func (c *virtualConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *virtualConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *virtualConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *virtualConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func multiplexConn(conn net.PacketConn) (stunConn net.PacketConn, appConn net.PacketConn) {
	return &virtualConn{conn}, &virtualConn{conn}
}
