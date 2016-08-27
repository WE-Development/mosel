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
	"encoding/json"
	"net/http"
	"bytes"
	"net/url"
	"log"
	"time"
	"bufio"
	"github.com/WE-Development/mosel/api"
)

type node struct {
	Name        string
	URL         url.URL
	scripts     []string

	handler     *nodeRespHandler
	scriptCache *scriptCache

	close       chan struct{}
}

func NewNode(name string, url url.URL, scripts []string, handler *nodeRespHandler, scriptCache *scriptCache) (*node, error) {

	node := &node{}
	node.Name = name
	node.URL = url
	node.scripts = scripts
	node.close = make(chan struct{})
	node.handler = handler
	node.scriptCache = scriptCache

	return node, nil
}

func (node *node) ListenStream() {
	Run: for {
		//provision scripts before connecting to stream
		if err := node.ProvisionScripts(); err != nil {
			log.Printf("Error while provisioning scripts: %s", err.Error())
		}

		url := node.URL.String() + "/stream"
		log.Printf("Connect to %s via %s", node.Name, url)
		resp, err := http.Get(url)

		if err != nil {
			log.Println(err)

			//todo make reconnection timeout configurable by moseld.conf
			time.Sleep(10 * time.Second)
			continue Run
		}

		for {
			select {
			case <-node.close:
				resp.Body.Close()
				break Run
			default:
				reader := bufio.NewReader(resp.Body)
				data, err := reader.ReadBytes('\n')

				if err != nil {
					//check weather we are dealing with a non-stream resource
					if err.Error() != "EOF" {
						log.Println(err)
					}

					resp.Body.Close()
					continue Run
				}

				var resp api.NodeResponse
				json.Unmarshal(data, &resp)
				node.handler.handleNodeResp(node.Name, resp)
			}
		}

	}
}

func (node *node) ProvisionScripts() error {
	for _, script := range node.scripts {
		var bytes []byte
		var err error

		if bytes, err = node.scriptCache.getScriptBytes(script); err != nil {
			return err
		}

		err = node.ProvisionScript(script, bytes)

		if err != nil {
			return err
		}
	}

	return nil
}

func (node *node) ProvisionScript(name string, b []byte) error {
	_, err := http.Post(node.URL.String() + "/script/" + name,
		"application/x-sh",
		bytes.NewBuffer(b))
	return err
}

