package ray

import (
	"context"
	"net/http"
)

var _ Protocol = (*HttpProtocol)(nil)

type HttpProtocol struct {
	trans Transport
}

func (p *HttpProtocol) OnConnect() {}

func (p *HttpProtocol) Read(ctx context.Context) (Packet, error) {
	return http.ReadRequest(p.trans.Reader())
}

func (p *HttpProtocol) Write(ctx context.Context, packet Packet) error {
	resp, ok := packet.(*http.Response)
	if !ok {
		return ErrInvalidPacket
	}
	return resp.Write(p.trans)
}

func (p *HttpProtocol) OnClose() {}
