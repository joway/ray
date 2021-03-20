package protocol

import (
	"context"
	"github/joway/ray"
)

type Protocol interface {
	OnConnect()
	Read(ctx context.Context) (ray.Packet, error)
	Write(ctx context.Context, packet ray.Packet) error
	OnClose()
}

type ProtocolFactory func(trans ray.Transport) Protocol
