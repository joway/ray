package protocol

import (
	"context"
	"github/joway/ray"
)

var _ Protocol = (*EchoProtocol)(nil)

type EchoProtocol struct {
	trans ray.Transport
	rbuf  []byte
}

func NewEchoProtocol(trans ray.Transport) Protocol {
	return &EchoProtocol{
		trans: trans,
		rbuf:  make([]byte, 1024),
	}
}

func (p *EchoProtocol) OnConnect() {
}

func (p *EchoProtocol) Read(ctx context.Context) (ray.Packet, error) {
	n, err := p.trans.Read(p.rbuf)
	if err != nil {
		return nil, err
	}
	return p.rbuf[:n], nil
}

func (p *EchoProtocol) Write(ctx context.Context, packet ray.Packet) error {
	pkt, ok := packet.([]byte)
	if !ok {
		return ray.ErrInvalidPacket
	}
	if _, err := p.trans.Write(pkt); err != nil {
		return err
	}
	return p.trans.Flush()
}

func (p *EchoProtocol) OnClose() {
}
