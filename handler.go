package sdk

import "github.com/it-chain/sdk/pb"

type TransactionHandler interface {
	Name() string
	Versions() []string
	Handle(request *pb.Request, cell *Cell) (*pb.Response, error)
}
