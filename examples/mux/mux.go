package main

import (
	"net"
	"time"
)

type virtualConn struct {
	readChan  <-chan readMsg
	writeChan chan<- writeMsg
}

type readMsg struct {
	p    []byte
	addr net.Addr
	err  error
}

type writeMsg struct {
	p    []byte
	addr net.Addr
}

func (c *virtualConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	// TODO: Add state to virtualConn to track when single packet hasn't been fully read
	msg := <-c.readChan
	n = copy(p, msg.p)
	addr, err = msg.addr, msg.err
	return
}

func (c *virtualConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	pCopy := make([]byte, len(p))
	n = copy(pCopy, p)

	msg := writeMsg{pCopy, addr}

	c.writeChan <- msg
	return
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
