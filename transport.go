package ray

import (
	"bufio"
	"io"
	"net"
)

type TransportFactory func(conn net.Conn) Transport

type Transport interface {
	io.ReadWriteCloser
	ReadBytes(delim byte) ([]byte, error)
	
	Reader() *bufio.Reader
	Writer() *bufio.Writer

	Peak(n int) ([]byte, error)
	Flush() error
}

type transport struct {
	*Socket
}

func NewTransport(conn net.Conn) Transport {
	return &transport{Socket: NewSocket(conn)}
}
