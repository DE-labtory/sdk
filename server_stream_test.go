/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package sdk_test

import (
	"context"
	"fmt"
	"io"
	"sync"
	"testing"
	"time"

	"strconv"

	"github.com/DE-labtory/sdk"
	"github.com/DE-labtory/sdk/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestServer_RunICode(t *testing.T) {
	server := sdk.NewServer(5002)
	wg := sync.WaitGroup{}
	wg.Add(1)
	server.SetHandler(func(request *pb.Request) *pb.Response {
		return &pb.Response{
			Uuid:  request.Uuid,
			Type:  "test",
			Data:  nil,
			Error: "",
		}
	})

	go func() {
		err := server.Listen(30)
		if err != nil {
			panic("error in listen server" + err.Error())
		}
	}()

	start := time.Now()
	stream, _ := DialMockClient("127.0.0.1:5002", func(response *pb.Response, e error) {
		fmt.Println(response.Uuid)
		if response.Uuid == "1999" {
			end := time.Now()
			assert.WithinDuration(t, end, start, 2*time.Second)
			wg.Done()
		}
	})
	for i := 0; i < 2000; i++ {
		stream.clientStream.Send(&pb.Request{
			Uuid:         strconv.Itoa(i),
			Type:         "invoke",
			FunctionName: "test_request",
			Args:         nil,
		})
	}

	wg.Wait()

}

func TestServer_Ping(t *testing.T) {

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
