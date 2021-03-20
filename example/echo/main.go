package main

import (
	"context"
	"errors"
	"github/joway/ray"
	"github/joway/ray/protocol"
	"log"
)

var (
	_ ray.ServerProcessor = (*EchoProcessor)(nil)
)

type EchoProcessor struct{}

func (p *EchoProcessor) Process(ctx context.Context, request ray.Packet) (ray.Packet, error) {
	req, ok := request.(*[]byte)
	if !ok {
		return nil, errors.New("cannot process invalid packet")
	}
	return req, nil
}

func main() {
	var (
		transFactory = ray.NewTransport
		protoFactory = protocol.NewEchoProtocol
		processor    = &EchoProcessor{}
	)
	server := ray.NewServer(transFactory, protoFactory, processor)

	if err := server.Serve(context.Background()); err != nil {
		log.Fatal(err)
	}
}
