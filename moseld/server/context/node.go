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

	"github.com/bluedevel/mosel/api"
)

// Represents a node and handles all communication with it.
type node struct {
	Name        string
	URL         url.URL
	scripts     []string

	user        string
	passwd      string

	handler     *nodeRespHandler
	scriptCache *scriptCache

	close       chan struct{}
}

func NewNode(name string, url url.URL, user  string, passwd string, scripts []string, handler *nodeRespHandler, scriptCache *scriptCache) (*node, error) {

	node := &node{}
	node.Name = name
	node.URL = url
	node.scripts = scripts
	node.user = user
	node.passwd = passwd
	node.close = make(chan struct{})
	node.handler = handler
	node.scriptCache = scriptCache

	return node, nil
}

// Connect and Listen to the node pushing infos.
// If the connection is lost, a reconnect is tried.
func (node *node) ListenStream() {
	Run: for {
		//provision scripts before connecting to stream
		if err := node.ProvisionScripts(); err != nil {
			log.Printf("Error while provisioning scripts: %s", err.Error())
		}

		// TRY
		resp, err := func() (*http.Response, error) {
			url := node.URL.String() + "/stream"
			log.Printf("Connect to %s via %s", node.Name, url)
			req, err := http.NewRequest("GET", url, nil)

			if err != nil {
				return nil, err
			}
			return node.doRequest(req)
		}()
		// CATCH
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

// Provision scripts with node scope to the node
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

// Provision a singe script to the node
func (node *node) ProvisionScript(name string, b []byte) error {
	err := func() (error) {
		url := node.URL.String() + "/script/" + name
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
		req.Header.Add("media-type", "application/x-sh")

		if err != nil {
			return err
		}
		_, err = node.doRequest(req)
		return err
	}()
	return err
}

// Do a request and setup basic auth if required for this node.
func (node *node) doRequest(req *http.Request) (*http.Response, error) {
	return func() (*http.Response, error) {
		if node.user != "" {
			req.SetBasicAuth(node.user, node.passwd)
		}

		client := http.Client{}
		resp, err := client.Do(req)

		return resp, err
	}()
}

