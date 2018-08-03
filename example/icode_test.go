package main

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	"sync"

	"github.com/it-chain/sdk"
	"github.com/it-chain/sdk/example/handler"
	"github.com/it-chain/sdk/logger"
	"github.com/it-chain/sdk/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestICode(t *testing.T) {
	logger.EnableFileLogger(true, "./icode.log")

	wg := sync.WaitGroup{}
	wg.Add(1)
	exHandler := &handler.HandlerExample{}
	ibox := sdk.NewIBox(5000)
	ibox.SetHandler(exHandler)
	go func() {
		err := ibox.On(30)
		if err != nil {
			panic("error in listen server" + err.Error())
		}
	}()
	start := time.Now()
	stream, _ := DialMockClient("127.0.0.1:5000", func(response *pb.Response, e error) {
		if e != nil {
			fmt.Println(e)
		}
		if response.Data != nil {
			fmt.Println("id : " + response.Uuid + ",res : " + string(response.Data))
		}
		if response.Uuid == "199" {
			end := time.Now()
			assert.WithinDuration(t, end, start, 2*time.Second)
			wg.Done()
		}
	})
	for i := 0; i < 200; i++ {
		stream.clientStream.Send(&pb.Request{
			Uuid:         strconv.Itoa(i),
			Type:         "test",
			FunctionName: "test_request",
			Args:         nil,
		})
	}

	wg.Wait()
}

type MockClient struct {
	clientStream pb.BistreamService_RunICodeClient
}

func DialMockClient(serverIp string, handler func(*pb.Response, error)) (*MockClient, context.CancelFunc) {

	dialContext, _ := context.WithTimeout(context.Background(), 3*time.Second)

	conn, err := grpc.DialContext(dialContext, serverIp, grpc.WithInsecure())
	if err != nil {
		panic("error in dial")
	}
	client := pb.NewBistreamServiceClient(conn)
	ctx, cf := context.WithCancel(context.Background())
	clientStream, err := client.RunICode(ctx)
	if err != nil {
		panic("error in run Icode" + err.Error())
	}
	mockClient := MockClient{
		clientStream: clientStream,
	}
	go func() {
		for {
			res, err := clientStream.Recv()
			if err == io.EOF {
				fmt.Println("io.EOF handle finish.")
				return
			}

			handler(res, err)

		}
	}()
	return &mockClient, cf
}
