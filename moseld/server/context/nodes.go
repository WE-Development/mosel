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
package context

import (
	"net/url"
	"net/http"
	"log"
	"bufio"
	"time"
	"encoding/json"
	"github.com/bluedevel/mosel/api"
	"errors"
)

type Node struct {
	Name  string
	URL   url.URL

	close chan struct{}
}

type nodeCache struct {
	nodes map[string]*Node
}

func NewNodeCache() (*nodeCache, error) {
	c := &nodeCache{}
	c.nodes = make(map[string]*Node)
	return c, nil
}

func (cache *nodeCache) Add(node *Node) {
	cache.nodes[node.Name] = node
	node.close = make(chan struct{})
	go func() {
		Connection: for {
			log.Printf("Connect to %s via %s", node.Name, node.URL.String())
			resp, err := http.Get(node.URL.String())

			if err != nil {
				log.Println(err)

				//todo make reconnection timeout configurable by moseld.conf
				time.Sleep(10 * time.Second)
				continue Connection
			}

			for {
				select {
				case <-node.close:
					resp.Body.Close()
					break Connection
				default:
					reader := bufio.NewReader(resp.Body)
					data, err := reader.ReadBytes('\n')

					if err != nil {
						//check weather we are dealing with a non-stream resource
						if err.Error() != "EOF" {
							log.Println(err)
						}

						resp.Body.Close()
						continue Connection
					}

					var resp api.NodeResponse
					json.Unmarshal(data, &resp)
					handleNodeResp(resp)
				}
			}

		}
	}()
}

func (cache *nodeCache) Get(name string) (*Node, error) {
	if val, ok := cache.nodes[name]; ok {
		return val, nil
	}

	return nil, errors.New("No node with name " + name + " registered")
}

func (cache *nodeCache) CloseNode(name string) error {
	node, err := cache.Get(name)

	if err != nil {
		return err
	}

	node.close <- struct{}{}
	return nil
}