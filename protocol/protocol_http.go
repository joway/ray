package protocol

import (
	"context"
	"github/joway/ray"
	"net/http"
)

var _ Protocol = (*HttpProtocol)(nil)

type HttpProtocol struct {
	trans ray.Transport
}

func (p *HttpProtocol) OnConnect() {}

func (p *HttpProtocol) Read(ctx context.Context) (ray.Packet, error) {
	return http.ReadRequest(p.trans.Reader())
}

func (p *HttpProtocol) Write(ctx context.Context, packet ray.Packet) error {
	resp, ok := packet.(*http.Response)
	if !ok {
		return ray.ErrInvalidPacket
	}
	return resp.Write(p.trans)
}

func (p *HttpProtocol) OnClose() {}
