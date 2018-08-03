package handler

import (
	"strconv"

	"errors"

	"github.com/it-chain/sdk"
	"github.com/it-chain/sdk/logger"
	"github.com/it-chain/sdk/pb"
)

type HandlerExample struct {
}

func (*HandlerExample) Name() string {
	return "sample"
}

func (*HandlerExample) Versions() []string {
	vers := make([]string, 0)
	vers = append(vers, "1.0")
	vers = append(vers, "1.2")
	return vers
}

func (*HandlerExample) Handle(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.Type {
	case "invoke":
		return handleInvoke(request, cell)
	case "query":
		return handleQuery(request, cell)
	default:
		logger.Fatal(nil, "unknown request type")
		err := errors.New("unknown request type")
		return responseError(request, err), err
	}
}
func handleQuery(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.FunctionName {
	case "getA":
		b, err := cell.GetData("A")
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, b), nil

	default:
		err := errors.New("unknown query method")
		return responseError(request, err), err
	}
}
func handleInvoke(request *pb.Request, cell *sdk.Cell) (*pb.Response, error) {
	switch request.FunctionName {
	case "initA":
		err := cell.PutData("A", []byte("0"))
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, nil), nil
	case "incA":
		data, err := cell.GetData("A")
		if err != nil {
			return responseError(request, err), err
		}
		if len(data) == 0 {
			err := errors.New("no data err")
			return responseError(request, err), err
		}
		strData := string(data)
		intData, err := strconv.Atoi(strData)
		if err != nil {
			return responseError(request, err), err
		}
		intData++
		changeData := strconv.Itoa(intData)
		err = cell.PutData("A", []byte(changeData))
		if err != nil {
			return responseError(request, err), err
		}
		return responseSuccess(request, nil), nil
	default:
		err := errors.New("unknown invoke method")
		return responseError(request, err), err
	}
}

func responseError(request *pb.Request, err error) *pb.Response {
	return &pb.Response{
		Uuid:   request.Uuid,
		Type:   request.Type,
		Result: false,
		Data:   nil,
		Error:  err.Error(),
	}
}

func responseSuccess(request *pb.Request, data []byte) *pb.Response {
	return &pb.Response{
		Uuid:   request.Uuid,
		Type:   request.Type,
		Result: true,
		Data:   data,
		Error:  "",
	}
}
