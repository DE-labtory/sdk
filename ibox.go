package sdk

import (
	"os"

	"github.com/it-chain/sdk/logger"
	"github.com/it-chain/sdk/pb"
)

type IBox struct {
	port    int
	handler TransactionHandler
}

func NewIBox(port int) *IBox {
	return &IBox{
		port:    port,
		handler: nil,
	}
}

func (i *IBox) SetHandler(handler TransactionHandler) {
	i.handler = handler
}

func (i *IBox) On(timeout int) error {
	if i.handler == nil {
		logger.Panic(nil, "handler is nil")
		os.Exit(1)
	}
	server := NewServer(i.port)
	cell := NewCell(i.handler.Name())

	serverHandler := func(request *pb.Request) (*pb.Response, error) {
		return i.handler.Handle(request, cell)
	}
	server.SetHandler(serverHandler)
	return server.Listen(timeout)
}
