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
	"errors"
	"log"
	"net/http"
	"time"
	"bufio"
	"github.com/WE-Development/mosel/api"
	"encoding/json"
	"bytes"
	"github.com/WE-Development/mosel/config"
)

type node struct {
	Name        string
	URL         url.URL
	scripts     []string

	handler     *nodeRespHandler
	scriptCache *scriptCache

	close       chan struct{}
}

func NewNode(name string, conf *moselconfig.NodeConfig, handler *nodeRespHandler, scriptCache *scriptCache) (*node, error) {

	node := &node{}
	node.Name = name
	node.close = make(chan struct{})
	node.handler = handler
	node.scriptCache = scriptCache

	return node, node.initialize(conf)
}

func (node *node) initialize(conf *moselconfig.NodeConfig) error {
	//get base url
	var url *url.URL
	var err error
	if url, err = url.Parse(conf.URL); err != nil {
		return err
	}

	var scripts []string

	//get configured
	if len(conf.Scripts) > 0 {
		scripts = conf.Scripts
	} else {
		scripts = node.scriptCache.GetScripts()
	}

	//exclude certain scripts
	for _, exclude := range conf.ScriptsExclude {
		for i, script := range scripts {
			if script == exclude {
				scripts = append(scripts[:i], scripts[i + 1:]...)
			}
		}
	}

	node.URL = *url
	node.scripts = scripts

	return nil;
}

func (node *node) ListenStream() {
	Connection: for {
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
				node.handler.handleNodeResp(node.Name, resp)
			}
		}

	}
}

func (node *node) ProvisionScripts() error {
	for _, script := range node.scripts {
		bytes, err := node.scriptCache.getScriptBytes(script)

		if err != nil {
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

type nodeCache struct {
	handler *nodeRespHandler
	scripts *scriptCache

	nodes   map[string]*node
}

func NewNodeCache(handler *nodeRespHandler, scripts *scriptCache) (*nodeCache, error) {
	c := &nodeCache{}
	c.handler = handler
	c.scripts = scripts
	c.nodes = make(map[string]*node)
	return c, nil
}

func (cache *nodeCache) Add(node *node) {
	cache.nodes[node.Name] = node

	//cache.ProvisionScripts(node.Name, cache.scripts.Scripts)
	go func() {
		node.ListenStream()
	}()
}

func (cache *nodeCache) Get(name string) (*node, error) {
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