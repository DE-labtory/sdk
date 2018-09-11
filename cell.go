/*
 * Copyright 2018 It-chain
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

package sdk

import (
	"github.com/it-chain/sdk/pb"
	"github.com/rs/xid"
)

type Cell struct {
	serverStream *Server
	icodeName    string
	requestMap   map[string]chan stateRequestResult
}

func NewCell(serverStream *Server, icodeName string) (*Cell, func(message *pb.Message)) {
	cellObj := &Cell{
		serverStream: serverStream,
		icodeName:    icodeName,
		requestMap:   make(map[string]chan stateRequestResult, 0),
	}
	return cellObj, cellObj.handleStateResult
}

func (c Cell) PutData(key string, value string) error {
	uuid := xid.New().String()
	mychan := make(chan stateRequestResult, 1)
	c.requestMap[uuid] = mychan
	err := c.serverStream.stream.Send(&pb.Message{
		Message: &pb.Message_WriteStateRequest{
			WriteStateRequest: &pb.WriteRequest{
				Uuid:      uuid,
				IcodeName: c.icodeName,
				Key:       key,
				Value:     value,
			},
		},
	})
	if err != nil {
		return err
	}
	for {
		select {
		case result := <-mychan:
			return result.err
		default:
		}
	}

}

func (c Cell) GetData(key string) ([]byte, error) {
	return value, err
}

func (c Cell) handleStateResult(message *pb.Message) {

}

type stateRequestResult struct {
	value string
	err   error
}
