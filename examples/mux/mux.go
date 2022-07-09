package main

import (
	"net"
	"time"
)

type logicalConn struct {
	conn net.PacketConn
}

func (c *logicalConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return c.conn.ReadFrom(p)
}

func (c *logicalConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	return c.conn.WriteTo(p, addr)
}
func (c *logicalConn) Close() error {
	return c.conn.Close()
}
func (c *logicalConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *logicalConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *logicalConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *logicalConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func multiplexConn(conn net.PacketConn) (stunConn net.PacketConn, appConn net.PacketConn) {
	return conn, conn
}
