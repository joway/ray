package ray

import (
	"bufio"
	"net"
	"sync"
)

var (
	_ Transport = (*Socket)(nil)

	socketReaderPool sync.Pool
	socketWriterPool sync.Pool
)

const (
	DefaultMaxBufferSize = 1024 * 1024 * 10 //10MB
	DefaultBufferSize    = 1024 * 8
)

type Socket struct {
	conn net.Conn
	rbuf *bufio.Reader
	wbuf *bufio.Writer
}

func NewSocket(conn net.Conn) *Socket {
	return &Socket{
		conn: conn,
		rbuf: acquireBufioReader(conn),
		wbuf: acquireBufioWriter(conn),
	}
}

func (s *Socket) Socket() *Socket {
	return s
}

func (s *Socket) Reader() *bufio.Reader {
	return s.rbuf
}

func (s *Socket) Writer() *bufio.Writer {
	return s.wbuf
}

func (s *Socket) Read(p []byte) (n int, err error) {
	return s.rbuf.Read(p)
}

func (s *Socket) ReadBytes(delim byte) ([]byte, error) {
	return s.rbuf.ReadBytes(delim)
}

func (s *Socket) Write(p []byte) (n int, err error) {
	return s.wbuf.Write(p)
}

func (s *Socket) Peak(n int) ([]byte, error) {
	return s.rbuf.Peek(n)
}

func (s *Socket) Flush() error {
	return s.wbuf.Flush()
}

func (s *Socket) Close() error {
	return s.conn.Close()
}

func acquireBufioReader(conn net.Conn) *bufio.Reader {
	rbuf := socketReaderPool.Get()
	if rbuf == nil {
		rbuf = bufio.NewReaderSize(conn, DefaultBufferSize)
	}
	rd := rbuf.(*bufio.Reader)
	rd.Reset(conn)
	return rd
}

func releaseBufioReader(buf *bufio.Reader) {
	if buf == nil {
		return
	}
	socketReaderPool.Put(buf)
}

func acquireBufioWriter(conn net.Conn) *bufio.Writer {
	wbuf := socketWriterPool.Get()
	if wbuf == nil {
		wbuf = bufio.NewWriterSize(conn, DefaultBufferSize)
	}
	wt := wbuf.(*bufio.Writer)
	wt.Reset(conn)
	return wt
}

func releaseBufioWriter(buf *bufio.Writer) {
	if buf == nil {
		return
	}
	socketWriterPool.Put(buf)
}
