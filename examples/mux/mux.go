package main

import (
	"net"
	"time"

	"github.com/cha-king/go-stun"
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

func NewVirtualConn(conn net.PacketConn) (stunConn net.PacketConn, appConn net.PacketConn) {
	// TODO: Sanity check buffer size
	readChanStun := make(chan readMsg, 1024)
	readChanApp := make(chan readMsg, 1024)
	writeChan := make(chan writeMsg, 1024)

	// TODO: Discard messages if buffers / channels are full

	// TODO: Contexts and things for gracefully handling goroutines

	// Reader goroutine
	go func() {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFrom(buf)
		// Safe to pass buf without copy, since consumer copies for us
		msg := readMsg{buf[:n], addr, err}

		if isStunMessage(buf[:n]) {
			readChanStun <- msg
		} else {
			readChanApp <- msg
		}
	}()

	// Writer goroutine
	go func() {
		msg := <-writeChan
		// TODO: Handle errors returned from WriteTo, maybe separate chan
		conn.WriteTo(msg.p, msg.addr)
	}()

	stunConn = &virtualConn{readChanStun, writeChan}
	appConn = &virtualConn{readChanApp, writeChan}
	return
}

func isStunMessage(p []byte) bool {
	// TODO: This should maybe be its own function. Probably useful for stun package.
	_, err := stun.DecodeHeader(p)
	return err == nil
}
