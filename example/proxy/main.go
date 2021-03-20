package main

import (
	"context"
	"errors"
	"github.com/joway/ray"
	"log"
	"net"
)

var (
	_ ray.Protocol        = (*TCPProtocol)(nil)
	_ ray.ServerProcessor = (*TCPProxyProcessor)(nil)
)

type TCPProtocol struct {
	trans ray.Transport
	rbuf  []byte
}

type TCPPacket []byte

func NewTCPProtocol(trans ray.Transport) ray.Protocol {
	return &TCPProtocol{
		trans: trans,
		rbuf:  make([]byte, 100),
	}
}

func (p *TCPProtocol) OnConnect() {
}

func (p *TCPProtocol) Read(ctx context.Context) (ray.Packet, error) {
	req := TCPPacket{}
	for {
		n, err := p.trans.Read(p.rbuf)
		if err != nil {
			return nil, err
		}
		req = append(req, p.rbuf[:n]...)
		if n < len(p.rbuf) {
			break
		}
	}
	return req, nil
}

func (p *TCPProtocol) Write(ctx context.Context, packet ray.Packet) error {
	data, ok := packet.(TCPPacket)
	if !ok {
		return errors.New("cannot write invalid packet")
	}
	if len(data) == 0 {
		return nil
	}
	if _, err := p.trans.Write(data); err != nil {
		return err
	}
	return p.trans.Flush()
}

func (p *TCPProtocol) OnClose() {
}

type TCPProxyProcessor struct {
	rbuf []byte
}

func NewTCPProxyProcessor() ray.ServerProcessor {
	return &TCPProxyProcessor{
		rbuf: make([]byte, 100),
	}
}

func (p *TCPProxyProcessor) Process(ctx context.Context, request ray.Packet) (ray.Packet, error) {
	req, ok := request.(TCPPacket)
	if !ok {
		return nil, errors.New("cannot process invalid packet")
	}

	conn, err := net.Dial("tcp", "httpbin.org:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(req)
	if err != nil {
		return nil, err
	}

	resp := TCPPacket{}
	for {
		n, err := conn.Read(p.rbuf)
		if err != nil {
			return nil, err
		}
		resp = append(resp, p.rbuf[:n]...)
		if n < len(p.rbuf) {
			break
		}
	}

	return resp, nil
}

func main() {
	var (
		transFactory = ray.NewTransport
		protoFactory = NewTCPProtocol
		processor    = NewTCPProxyProcessor()
	)
	server := ray.NewServer(transFactory, protoFactory, processor)

	if err := server.Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
