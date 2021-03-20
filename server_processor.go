package ray

import (
	"context"
)

type ServerProcessor interface {
	Process(ctx context.Context, request Packet) (Packet, error)
}
