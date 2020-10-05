package http_server

import (
	"errors"
	"net"
)

type listenerDialer struct {
	// server conns to accept
	conns chan net.Conn
}

func NewListenerDialer() *listenerDialer {
	ld := &listenerDialer{
		conns: make(chan net.Conn),
	}
	return ld
}

// net.Listener interface implementation
func (ld *listenerDialer) Accept() (net.Conn, error) {
	conn, ok := <-ld.conns
	if !ok {
		return nil, errors.New("listenerDialer is closed")
	}
	return conn, nil
}

// net.Listener interface implementation
func (ld *listenerDialer) Close() error {
	close(ld.conns)
	return nil
}

// net.Listener interface implementation
func (ld *listenerDialer) Addr() net.Addr {
	// return arbitrary fake addr.
	return &net.UnixAddr{
		Name: "listenerDialer",
		Net:  "fake",
	}
}

// Transport.Dial implementation
func (ld *listenerDialer) Dial(network, addr string) (net.Conn, error) {
	cConn, sConn := net.Pipe()
	ld.conns <- sConn
	return cConn, nil
}
