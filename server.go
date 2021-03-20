package ray

import (
	"context"
	"errors"
	"fmt"
	"github/joway/ray/protocol"
	"io"
	"net"
)

const (
	DefaultServerPort = 9090
)

var (
	ErrServerTerminated = errors.New("server has been terminated")
)

type Server struct {
	ip   string
	port int

	transFactory TransportFactory
	protoFactory protocol.ProtocolFactory
	processor    ServerProcessor

	logger Logger

	terminated chan struct{}
}

func defaultServer() *Server {
	return &Server{
		port:       DefaultServerPort,
		logger:     NewDefaultLogger(),
		terminated: make(chan struct{}),
	}
}

func NewServer(
	transportFactory TransportFactory, protocolFactory protocol.ProtocolFactory,
	processor ServerProcessor,
	options ...ServerOption,
) *Server {
	server := defaultServer()
	for _, optFunc := range options {
		optFunc(server)
	}
	server.transFactory = transportFactory
	server.protoFactory = protocolFactory
	server.processor = processor
	return server
}

func (s *Server) Serve(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.ip, s.port)
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					s.logger.Errorf("process panic: %v", err)
				}
			}()

			if err := s.process(ctx, conn); err != nil {
				switch err {
				case io.EOF:
					break
				default:
					s.logger.Warnf("server process failed: %v", err)
				}
			}
		}()
	}
}

func (s *Server) process(ctx context.Context, conn net.Conn) error {
	trans := s.transFactory(conn)
	proto := s.protoFactory(trans)

	defer func() {
		proto.OnClose()

		if err := conn.Close(); err != nil {
			s.logger.Warnf("server process close conn failed: %v", err)
		}
	}()

	proto.OnConnect()
	for {
		select {
		case <-s.terminated:
			return ErrServerTerminated
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		request, err := proto.Read(ctx)
		if err != nil {
			return err
		}

		response, err := s.processor.Process(ctx, request)
		if err != nil {
			return err
		}

		if err := proto.Write(ctx, response); err != nil {
			return err
		}
	}
}

func (s *Server) Stop() error {
	select {
	case <-s.terminated:
	default:
		close(s.terminated)
	}
	return nil
}
