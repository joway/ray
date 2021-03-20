package ray

import (
	"context"
)

type Protocol interface {
	OnConnect()
	Read(ctx context.Context) (Packet, error)
	Write(ctx context.Context, packet Packet) error
	OnClose()
}

type ProtocolFactory func(trans Transport) Protocol
