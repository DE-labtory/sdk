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

import "github.com/it-chain/sdk/pb"

type RequestHandler interface {
	// this function returns icode name
	Name() string
	// this function returns icode versions that can be handled
	Versions() []string
	// this function returns icode names that can read icode data.
	// for example, let your icode name is A. if ReaderList() return {"B"}, then
	// icode name B can access A icode data for read permission
	ReaderList() []string
	// this function returns icode names that can write icode data.
	// for example, let your icode name is A. if ReaderList() return {"B"}, then
	// icode name B can access A icode data for write permission
	WriterList() []string
	// this function handle request and return response.
	Handle(request *pb.RunICodeRequest, cell *Cell) *pb.RunICodeResponse
}
