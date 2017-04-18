/*
 * Copyright 2016 Robin Engel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package api

/*
Represents the state of a node at a certain time.
Diagram -> Graph -> Value
*/
type NodeInfo map[string]map[string]string

//A wrapper for sending data over the stream from the nodes to server
type NodeResponse struct {
	moselResponse
	NodeInfo NodeInfo `json:"nodeÌnfo"`
}

func NewNodeResponse() NodeResponse {
	return NodeResponse{
		moselResponse: newMoselResponse(),
		NodeInfo:      make(NodeInfo),
	}
}
